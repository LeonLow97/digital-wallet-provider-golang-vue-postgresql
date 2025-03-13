package jsonutil_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LeonLow97/go-clean-architecture/utils/jsonutil"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestReadJSONBody(t *testing.T) {
	type TestStruct struct {
		Field string
	}

	tests := []struct {
		name    string
		body    string
		want    TestStruct
		wantErr error
	}{
		{
			name: "ValidJSON",
			body: `{"Field": "value"}`,
			want: TestStruct{
				Field: "value",
			},
		},
		{
			name:    "InvalidJSON",
			body:    `{"Field": "value}`,
			wantErr: errors.New("unexpected EOF"),
		},
		{
			name:    "EmptyBody",
			body:    "",
			wantErr: errors.New("EOF"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/", strings.NewReader(tt.body))
			require.NoError(t, err)

			var dst TestStruct
			err = jsonutil.ReadJSONBody(nil, req, &dst)

			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.want, dst)
		})
	}
}

func TestReadURLParamsInt(t *testing.T) {
	t.Run("ValidParam", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/123", nil)
		w := httptest.NewRecorder()
		req = mux.SetURLVars(req, map[string]string{"id": "123"})

		id, err := jsonutil.ReadURLParamsInt(w, req, "id")
		require.NoError(t, err)
		require.Equal(t, 123, id)
	})

	t.Run("MissingParam", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		id, err := jsonutil.ReadURLParamsInt(w, req, "id")
		require.Error(t, err)
		require.Equal(t, 0, id)
		require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("InvalidParam", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/abc", nil)
		w := httptest.NewRecorder()
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})

		id, err := jsonutil.ReadURLParamsInt(w, req, "id")
		require.Error(t, err)
		require.Equal(t, 0, id)
		require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})
}

func TestReadURLParamsString(t *testing.T) {
	t.Run("ValidParam", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/value", nil)
		w := httptest.NewRecorder()
		req = mux.SetURLVars(req, map[string]string{"key": "value"})

		value, err := jsonutil.ReadURLParamsString(w, req, "key")
		require.NoError(t, err)
		require.Equal(t, "value", value)
	})

	t.Run("MissingParam", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		value, err := jsonutil.ReadURLParamsString(w, req, "key")
		require.Error(t, err)
		require.Equal(t, "", value)
		require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})
}

func TestReadQueryParams(t *testing.T) {
	t.Run("ValidQueryParams", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?key=value", nil)
		var params struct {
			Key string `form:"key"`
		}

		err := jsonutil.ReadQueryParams(&params, req)
		require.NoError(t, err)
		require.Equal(t, "value", params.Key)
	})

	t.Run("InvalidQueryParams", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?key=value", nil)
		var params struct {
			Key int `form:"key"` // incompatible type to force an error
		}

		err := jsonutil.ReadQueryParams(&params, req)
		require.Error(t, err)
		require.Equal(t, err.Error(), "Field Namespace:key ERROR:Invalid Integer Value 'value' Type 'int' Namespace 'key'")
	})
}
