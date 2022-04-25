package router

import (
	"api/src/controllers"
	"net/http"
)

var routerUser = []Route{
	{
		URI:                    "/users",
		Methods:                http.MethodPost,
		Function:               controllers.CreateUser,
		RequiresAuthentication: false, // Temporary test
	},
	{
		URI:                    "/users",
		Methods:                http.MethodGet,
		Function:               controllers.SearchUsers,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/users/{userId}",
		Methods:                http.MethodGet,
		Function:               controllers.SearchUser,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/users/{userId}",
		Methods:                http.MethodPut,
		Function:               controllers.UpdateUser,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/users/{userId}",
		Methods:                http.MethodDelete,
		Function:               controllers.DeleteUser,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/users/{userId}/seguir",
		Methods:                http.MethodPost,
		Function:               controllers.FollowUser,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/users/{userId}/parar-de-seguir",
		Methods:                http.MethodPost,
		Function:               controllers.UnfollowollowUser,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/users/{userId}/followers",
		Methods:                http.MethodPost,
		Function:               controllers.SearchFollowers,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/users/{userId}/following",
		Methods:                http.MethodPost,
		Function:               controllers.SearchFollowing,
		RequiresAuthentication: true,
	},
	{
		URI:                    "/users/{userId}/update-password",
		Methods:                http.MethodPost,
		Function:               controllers.UpdatePassword,
		RequiresAuthentication: true,
	},
}
