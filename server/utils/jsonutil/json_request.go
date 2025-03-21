package jsonutil

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	apiErr "github.com/LeonLow97/go-clean-architecture/exception/response"
	"github.com/go-playground/form"
	"github.com/gorilla/mux"
)

// ReadJSONBody decodes the JSON body from an HTTP request into the provided struct.
func ReadJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		log.Println("error decoding request body:", err)
		return err
	}

	return nil
}

// ReadURLParamsInt extracts an integer parameter from the URL.
func ReadURLParamsInt(w http.ResponseWriter, r *http.Request, paramName string) (int, error) {
	vars := mux.Vars(r)
	paramValueString, ok := vars[paramName]
	if !ok {
		log.Printf("missing required param: %s\n", paramName)
		ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return 0, errors.New("missing required param")
	}

	paramValue, err := strconv.Atoi(paramValueString)
	if err != nil {
		log.Printf("invalid value for param: %s with error: %v\n", paramName, err)
		ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return 0, errors.New("invalid value for param")
	}

	return paramValue, nil
}

// ReadURLParamsString extracts a string parameter from the URL.
func ReadURLParamsString(w http.ResponseWriter, r *http.Request, paramName string) (string, error) {
	vars := mux.Vars(r)
	paramValue, ok := vars[paramName]
	if !ok {
		log.Printf("missing required param: %s\n", paramName)
		ErrorJSON(w, apiErr.ErrBadRequest, http.StatusBadRequest)
		return "", errors.New("missing required param")
	}
	return paramValue, nil
}

// ReadQueryParams decodes query parameters into a struct
func ReadQueryParams(dest interface{}, r *http.Request) error {
	if err := form.NewDecoder().Decode(dest, r.URL.Query()); err != nil {
		log.Printf("failed to read query params with error: %v\n", err)
		return err
	}
	return nil
}
