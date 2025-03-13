package jsonutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/LeonLow97/go-clean-architecture/utils/constants/headers"
	"github.com/LeonLow97/go-clean-architecture/utils/pagination"
)

// errorResponse is the type used for sending error json response
type errorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// ReadJSON tries to read the body of a request and converts from json into a go data variable
func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1024 * 1024 // 1 MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(data)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "http: request body too large"):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
		case errors.As(err, &invalidUnmarshalError):
			return fmt.Errorf("error unmarshaling JSON: %s", err.Error())
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{}) // decode more JSON from that file
	if err != io.EOF {
		return errors.New("body must contain only one JSON value")
	}

	return nil
}

// WriteJSON takes a response status code and arbitrary data and write json to the client
func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

// WriteNoContent takes a response status code with no content
func WriteNoContent(w http.ResponseWriter, status int, headers ...http.Header) {
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
}

// ErrorJSON takes an error message, and optionally a status code, generates and sends an error response json
func ErrorJSON(w http.ResponseWriter, errMessage string, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload errorResponse
	payload.Status = statusCode
	payload.Message = errMessage

	return WriteJSON(w, statusCode, payload)
}

// SetPaginatorHeaders sets the paginator details in the response headers
func SetPaginatorHeaders(w http.ResponseWriter, paginator *pagination.Paginator) {
	if paginator.Page > paginator.TotalPages() {
		paginator.Page = paginator.TotalPages()
	}

	w.Header().Set(headers.XTotal, strconv.FormatInt(paginator.TotalRecords, 10))
	w.Header().Set(headers.XTotalPages, strconv.FormatInt(paginator.TotalPages(), 10))
	w.Header().Set(headers.XPage, strconv.FormatInt(paginator.Page, 10))
	w.Header().Set(headers.XPageSize, strconv.FormatInt(paginator.PageSize, 10))
	w.Header().Set(headers.XHasNextPage, strconv.FormatBool(paginator.HasNextPage()))
	w.Header().Set(headers.XHasPreviousPage, strconv.FormatBool(paginator.HasPreviousPage()))
}
