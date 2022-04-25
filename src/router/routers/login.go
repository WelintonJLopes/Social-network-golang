package router

import (
	"api/src/controllers"
	"net/http"
)

var routerLogin = Route{
	URI:                    "/login",
	Methods:                http.MethodPost,
	Function:               controllers.Login,
	RequiresAuthentication: false,
}
