package router

import (
	router "api/src/router/routers"

	"github.com/gorilla/mux"
)

func Generate() *mux.Router {
	r := mux.NewRouter()
	return router.Configure(r)
}
