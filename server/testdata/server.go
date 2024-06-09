package testdata

import (
	"net/http/httptest"

	"github.com/gorilla/mux"
)

func Server(middlewares ...mux.MiddlewareFunc) (*httptest.Server, *mux.Router) {
	router := mux.NewRouter()
	for _, mw := range middlewares {
		router.Use(mw)
	}

	server := httptest.NewServer(router)
	return server, router
}
