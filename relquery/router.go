package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

// NewRouter forms a new mux router, see https://github.com/gorilla/mux
func NewRouter() *mux.Router {

	// create a basic router
	router := mux.NewRouter().StrictSlash(true)

	// assign the handlers to run when endpoints are called
	for _, route := range routes {

		var handler http.Handler
		handler = route.HandlerFunc

		// wrap all current routes in the logger decorator to log out requests
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}
