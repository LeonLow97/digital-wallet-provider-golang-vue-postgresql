package jsonutil_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LeonLow97/go-clean-architecture/utils/constants/headers"
	"github.com/LeonLow97/go-clean-architecture/utils/jsonutil"
	"github.com/LeonLow97/go-clean-architecture/utils/pagination"
	"github.com/stretchr/testify/require"
)

// Mock struct for testing
type TestStruct struct {
	Name string `json:"name"`
}

func TestReadJSON(t *testing.T) {
	t.Run("ValidJSON", func(t *testing.T) {
		body := `{"name":"test"}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		w := httptest.NewRecorder()

		var data TestStruct
		err := jsonutil.ReadJSON(w, req, &data)
		require.NoError(t, err)
		require.Equal(t, "test", data.Name)
	})

	t.Run("MalformedJSON", func(t *testing.T) {
		body := `{"name":"test"`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		w := httptest.NewRecorder()

		var data TestStruct
		err := jsonutil.ReadJSON(w, req, &data)
		require.Error(t, err)
		require.Contains(t, err.Error(), "body contains badly formed JSON")
	})

	t.Run("IncorrectJSONType", func(t *testing.T) {
		body := `{"name":123}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		w := httptest.NewRecorder()

		var data TestStruct
		err := jsonutil.ReadJSON(w, req, &data)
		require.Error(t, err)
		require.Contains(t, err.Error(), "body contains incorrect JSON type for field \"name\"")
	})

	t.Run("UnknownField", func(t *testing.T) {
		body := `{"unknown":"field"}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		w := httptest.NewRecorder()

		var data TestStruct
		err := jsonutil.ReadJSON(w, req, &data)
		require.Error(t, err)
		require.Contains(t, err.Error(), "json: unknown field")
	})

	t.Run("EmptyBody", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
		w := httptest.NewRecorder()

		var data TestStruct
		err := jsonutil.ReadJSON(w, req, &data)
		require.Error(t, err)
		require.Contains(t, err.Error(), "body must not be empty")
	})

	// t.Run("too large body", func(t *testing.T) {
	// 	largeBody := strings.Repeat("a", 1024*1024+1)
	// 	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(largeBody))
	// 	w := httptest.NewRecorder()

	// 	var data TestStruct
	// 	err := jsonutil.ReadJSON(w, req, &data)
	// 	require.Error(t, err)
	// 	require.Contains(t, err.Error(), "body must not be larger than 1048576 bytes")
	// })

	t.Run("multiple JSON values", func(t *testing.T) {
		body := `{"name":"test"}{"name":"test2"}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		w := httptest.NewRecorder()

		var data TestStruct
		err := jsonutil.ReadJSON(w, req, &data)
		require.Error(t, err)
		require.Contains(t, err.Error(), "body must contain only one JSON value")
	})
}

func TestWriteJSON(t *testing.T) {
	t.Run("ValidResponse", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := TestStruct{Name: "test"}

		err := jsonutil.WriteJSON(w, http.StatusOK, data)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var respData TestStruct
		err = json.NewDecoder(w.Body).Decode(&respData)
		require.NoError(t, err)
		require.Equal(t, data, respData)
	})

	t.Run("WithHeaders", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := TestStruct{Name: "test"}
		headers := http.Header{"X-Custom-Header": []string{"value"}}

		err := jsonutil.WriteJSON(w, http.StatusOK, data, headers)
		require.NoError(t, err)
		require.Equal(t, "value", w.Header().Get("X-Custom-Header"))
	})
}

func TestWriteNoContent(t *testing.T) {
	t.Run("NoContentResponse", func(t *testing.T) {
		w := httptest.NewRecorder()
		jsonutil.WriteNoContent(w, http.StatusNoContent)
		require.Equal(t, http.StatusNoContent, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	})

	t.Run("WithHeaders", func(t *testing.T) {
		w := httptest.NewRecorder()
		headers := http.Header{"X-Custom-Header": []string{"value"}}
		jsonutil.WriteNoContent(w, http.StatusNoContent, headers)
		require.Equal(t, http.StatusNoContent, w.Code)
		require.Equal(t, "value", w.Header().Get("X-Custom-Header"))
	})
}

// errorResponse is the type used for sending error json response
type errorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func TestErrorJSON(t *testing.T) {
	t.Run("DefaultStatusCode", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := jsonutil.ErrorJSON(w, "error message")
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var resp errorResponse
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, errorResponse{Status: http.StatusBadRequest, Message: "error message"}, resp)
	})

	t.Run("CustomStatusCode", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := jsonutil.ErrorJSON(w, "error message", http.StatusInternalServerError)
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var resp errorResponse
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, errorResponse{Status: http.StatusInternalServerError, Message: "error message"}, resp)
	})
}

func TestSetPaginatorHeaders(t *testing.T) {
	tests := []struct {
		name            string
		paginator       *pagination.Paginator
		expectedHeaders map[string]string
	}{
		{
			name: "ValidPaginator",
			paginator: &pagination.Paginator{
				Page:         1,
				PageSize:     10,
				TotalRecords: 100,
			},
			expectedHeaders: map[string]string{
				headers.XTotal:           "100",
				headers.XTotalPages:      "10",
				headers.XPage:            "1",
				headers.XPageSize:        "10",
				headers.XHasNextPage:     "true",
				headers.XHasPreviousPage: "false",
			},
		},
		{
			name: "PageGreaterThanTotalPages",
			paginator: &pagination.Paginator{
				Page:         11,
				PageSize:     10,
				TotalRecords: 100,
			},
			expectedHeaders: map[string]string{
				headers.XTotal:           "100",
				headers.XTotalPages:      "10",
				headers.XPage:            "10",
				headers.XPageSize:        "10",
				headers.XHasNextPage:     "false",
				headers.XHasPreviousPage: "true",
			},
		},
		{
			name: "ZeroTotalRecords",
			paginator: &pagination.Paginator{
				Page:         1,
				PageSize:     10,
				TotalRecords: 0,
			},
			expectedHeaders: map[string]string{
				headers.XTotal:           "0",
				headers.XTotalPages:      "0",
				headers.XPage:            "0",
				headers.XPageSize:        "10",
				headers.XHasNextPage:     "false",
				headers.XHasPreviousPage: "false",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			jsonutil.SetPaginatorHeaders(w, tt.paginator)

			for header, expectedValue := range tt.expectedHeaders {
				actualValue := w.Header().Get(header)
				require.Equal(t, expectedValue, actualValue, fmt.Sprintf("Header: %s", header))
			}
		})
	}
}
