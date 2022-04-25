package router

import (
	"api/src/controllers"
	"net/http"
)

var routerPublication = []Route{
	{
		URI:                    "/publications",
		Methods:                http.MethodPost,
		Function:               controllers.CreatePublication,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/publications",
		Methods:                http.MethodGet,
		Function:               controllers.SearchPublications,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/publications/{publicationId}",
		Methods:                http.MethodPost,
		Function:               controllers.SearchPublication,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/publications/{publicationId}",
		Methods:                http.MethodPut,
		Function:               controllers.UpdatePublication,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/publications/{publicationId}",
		Methods:                http.MethodDelete,
		Function:               controllers.DeletePublication,
		RequiresAuthentication: true,
	},
	{
		URI:                    "users/{usersId}/publications",
		Methods:                http.MethodGet,
		Function:               controllers.SearchPublicationsUser,
		RequiresAuthentication: true,
	},
	{
		URI:                    "publications/{publicationId}/like",
		Methods:                http.MethodPost,
		Function:               controllers.LikePublication,
		RequiresAuthentication: true,
	},
	{
		URI:                    "publications/{publicationId}/dislike",
		Methods:                http.MethodPost,
		Function:               controllers.DislikePublication,
		RequiresAuthentication: true,
	},
}
