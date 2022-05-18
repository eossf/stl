package main

import "github.com/julienschmidt/httprouter"

/*
Define all the routes here.
A new Route entry passed to the routes slice will be automatically
translated to a handler with the NewRouter() function
*/
type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

type Routes []Route

func AllRoutes() Routes {
	routes := Routes{
		Route{"Index", "GET", "/", Index},
		Route{"TrackIndex", "GET", "/tracks", TrackIndex},
		Route{"TrackShow", "GET", "/tracks/:id", TrackShow},
		Route{"TrackCreate", "POST", "/tracks", TrackCreate},
		Route{"TrackUpdate", "PUT", "/tracks/:id", TrackUpdate},
	}
	return routes
}
