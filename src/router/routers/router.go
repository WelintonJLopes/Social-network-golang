package router

import (
	"api/src/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	URI                    string
	Methods                string
	Function               func(http.ResponseWriter, *http.Request)
	RequiresAuthentication bool
}

func Configure(r *mux.Router) *mux.Router {
	router := routerUser
	router = append(router, routerLogin)
	router = append(router, routerPublication...)

	for _, route := range router {

		if route.RequiresAuthentication {
			r.HandleFunc(route.URI, middlewares.Logger(middlewares.Authenticate(route.Function))).Methods(route.Methods)
		} else {
			r.HandleFunc(route.URI, middlewares.Logger(route.Function)).Methods(route.Methods)
		}

	}

	return r
}
