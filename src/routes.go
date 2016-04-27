package main

import (
	"net/http"

	"github.com/freehaha/token-auth"
	"github.com/freehaha/token-auth/memory"
	"github.com/gorilla/mux"
)

// Route type is used to define a route of the API
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes type is an array of Route
type Routes []Route

// NewRouter is the constructeur of the Router
// It will create every routes from the routes variable just above
func NewRouter() *mux.Router {
	tokenAuthUser := tauth.NewTokenAuth(nil, nil, memStoreUser, nil)
	tokenAuthSuperUser := tauth.NewTokenAuth(nil, nil, memStoreSuperUser, nil)

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range publicRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(tokenAuthUser.HandleFunc(route.HandlerFunc))
	}

	for _, route := range superRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(tokenAuthSuperUser.HandleFunc(route.HandlerFunc))
	}

	return router
}

var memStoreUser = memstore.New("user")
var memStoreSuperUser = memstore.New("superUser")

var publicRoutes = Routes{
	Route{"LogAssociation", "POST", "/login/association", LogAssociationController},
}

var superRoutes = Routes{
	Route{"AddAssociation", "POST", "/association", AddAssociationController},
	Route{"DeleteAssociation", "DELETE", "/association/{id}", DeleteAssociationController},
	Route{"CreateUserForAssociation", "POST", "/association/{id}/user", CreateUserForAssociationController},
}

var routes = Routes{
	Route{"Index", "GET", "/", Index},

	//ASSOCIATIONS
	Route{"GetAssociation", "GET", "/association", GetAllAssociationsController},
	Route{"GetAssociation", "GET", "/association/{id}", GetAssociationController},
	Route{"UpdateAssociation", "PUT", "/association/{id}", UpdateAssociationController},
	Route{"ImageAssociation", "POST", "/association/{id}/image", AddImageAssociationController},

	//EVENTS
	Route{"GetFutureEvents", "GET", "/event", GetFutureEventsController},
	Route{"GetEvent", "GET", "/event/{id}", GetEventController},
	Route{"AddEvent", "POST", "/event", AddEventController},
	Route{"UpdateEvent", "PUT", "/event/{id}", UpdateEventController},
	Route{"DeleteEvent", "DELETE", "/event/{id}", DeleteEventController},
	Route{"ImageEvent", "POST", "/event/{id}/image", AddImageEventController},
	Route{"AddParticipant", "POST", "/event/{id}/participant/{userID}", AddParticipantController},
	Route{"RemoveParticipant", "DELETE", "/event/{id}/participant/{userID}", RemoveParticipantController},

	//POSTS
	Route{"GetPost", "GET", "/post/{id}", GetPostController},
	Route{"GetLastestPost", "GET", "/post", GetLastestPostsController},
	Route{"AddPost", "POST", "/post", AddPostController},
	Route{"UpdatePost", "PUT", "/post/{id}", UpdatePostController},
	Route{"DeletePost", "DELETE", "/post/{id}", DeletePostController},
	Route{"ImagePost", "POST", "/post/{id}/image", AddImagePostController},
	Route{"LikePost", "POST", "/post/{id}/like/{userID}", LikePostController},
	Route{"DislikePost", "DELETE", "/post/{id}/like/{userID}", DislikePostController},
	Route{"CommentPost", "POST", "/post/{id}/comment", CommentPostController},
	Route{"UncommentPost", "DELETE", "/post/{id}/comment/{commentID}", UncommentPostController},

	//USER
	Route{"GetUser", "GET", "/user/{id}", GetUserController},
	Route{"AddUser", "POST", "/user", AddUserController},
	Route{"UpdateUser", "PUT", "/user/{id}", UpdateUserController},
	Route{"DeleteUser", "DELETE", "/user/{id}", DeleteUserController},
	Route{"ImageUser", "POST", "/user/{id}/image", AddImageUserController},
}
