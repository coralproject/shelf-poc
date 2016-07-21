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
		"GetUser",
		"GET",
		"/user",
		GetUser,
	},
	Route{
		"GetComment",
		"GET",
		"/comment",
		GetComment,
	},
	Route{
		"GraphQuerySingle",
		"GET",
		"/graph/singleasset",
		GraphQuerySingle,
	},
	Route{
		"MongoQuerySingle",
		"GET",
		"/mongo/singleasset",
		MongoQuerySingle,
	},
	Route{
		"GraphQueryUserAssets",
		"GET",
		"/graph/userassets",
		GraphQueryUserAssets,
	},
	Route{
		"MongoQueryUserAssets",
		"GET",
		"/mongo/userassets",
		MongoQueryUserAssets,
	},
	Route{
		"GraphQueryUserComments",
		"GET",
		"/graph/usercomments",
		GraphQueryUserComments,
	},
	Route{
		"MongoQueryUserComments",
		"GET",
		"/mongo/usercomments",
		MongoQueryUserComments,
	},
	Route{
		"GraphQueryParComments",
		"GET",
		"/graph/parentcomments",
		GraphQueryParComments,
	},
	Route{
		"MongoQueryParComments",
		"GET",
		"/mongo/parentcomments",
		MongoQueryParComments,
	},
}
