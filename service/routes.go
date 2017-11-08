package service

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = Routes{
	Route{
		"getDomain",
		"GET",
		"/v1/domain/{fqdn}",
		getDomain,
	},
	Route{
		"createDomain",
		"POST",
		"/v1/domain",
		createDomain,
	},
	Route{
		"updateDomain",
		"PUT",
		"/v1/domain/{fqdn}",
		updateDomain,
	},
	Route{
		"deleteDomain",
		"DELETE",
		"/v1/domain/{fqdn}",
		deleteDomain,
	},
	Route{
		"renewDomain",
		"PUT",
		"/v1/domain/{fqdn}",
		renewDomain,
	},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(apiHandler(route.HandlerFunc))
	}

	return router
}