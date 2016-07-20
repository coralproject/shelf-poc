package main

import "net/http"

// Route is used to pass information about a particular route.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes is used to pass information about multiple routes.
type Routes []Route

// routes are the routes made available by this server.
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
