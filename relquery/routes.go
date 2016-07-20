package main

import "net/http"

// Route - used to pass information about a particular route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes - used to pass information about multiple routes
type Routes []Route

var routes = Routes{
	Route{
		"GetAsset",
		"GET",
		"/asset",
		GetAsset,
	},
	Route{
		"GraphQuery",
		"GET",
		"/graph",
		GraphQuery,
	},
	Route{
		"MongoQuery",
		"GET",
		"/mongo",
		MongoQuery,
	},
}
