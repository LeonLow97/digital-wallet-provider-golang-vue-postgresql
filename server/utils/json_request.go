package utils

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ReadParamsInt extracts an integer parameter from the URL.
func ReadParamsInt(r *http.Request, paramName string) (int, error) {
	vars := mux.Vars(r)
	paramValueString, ok := vars[paramName]
	if !ok {
		log.Printf("missing required param: %s\n", paramName)
		return 0, errors.New("missing required param")
	}

	paramValue, err := strconv.Atoi(paramValueString)
	if err != nil {
		log.Printf("invalid value for param: %s with error: %v\n", paramName, err)
		return 0, errors.New("invalid value for param")
	}

	return paramValue, nil
}
