package main

import (
	"net/http"

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
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

var routes = Routes{
	Route{"Index", "GET", "/", Index},

	//ASSOCIATIONS
	Route{"GetAssociation", "GET", "/association", GetAllAssociationsController},
	Route{"GetAssociation", "GET", "/association/{id}", GetAssociationController},
	Route{"AddAssociation", "POST", "/association", AddAssociationController},
	Route{"UpdateAssociation", "PUT", "/association/{id}", UpdateAssociationController},
	Route{"DeleteAssociation", "DELETE", "/association/{id}", DeleteAssociationController},

	//EVENTS
	Route{"GetFutureEvents", "GET", "/event", GetFutureEventsController},
	Route{"GetEvent", "GET", "/event/{id}", GetEventController},
	Route{"AddEvent", "POST", "/event", AddEventController},
	Route{"UpdateEvent", "PUT", "/event/{id}", UpdateEventController},
	Route{"DeleteEvent", "DELETE", "/event/{id}", DeleteEventController},
	Route{"AddParticipant", "POST", "/event/{id}/participant/{userID}", AddParticipantController},
	Route{"RemoveParticipant", "DELETE", "/event/{id}/participant/{userID}", RemoveParticipantController},

	//POSTS
	Route{"GetPost", "GET", "/post/{id}", GetPostController},
	Route{"GetLastestPost", "GET", "/post", GetLastestPostsController},
	Route{"AddPost", "POST", "/post", AddPostController},
	Route{"UpdatePost", "PUT", "/post/{id}", UpdatePostController},
	Route{"DeletePost", "DELETE", "/post/{id}", DeletePostController},
	Route{"LikePost", "POST", "/post/{id}/like/{userID}", LikePostController},
	Route{"DislikePost", "DELETE", "/post/{id}/like/{userID}", DislikePostController},
	Route{"CommentPost", "POST", "/post/{id}/comment", CommentPostController},
	Route{"UncommentPost", "DELETE", "/post/{id}/comment/{commentID}", UncommentPostController},

	//USER
	Route{"GetUser", "GET", "/user/{id}", GetUserController},
	Route{"AddUser", "POST", "/user", AddUserController},
	Route{"UpdateUser", "PUT", "/user/{id}", UpdateUserController},
	Route{"DeleteUser", "DELETE", "/user/{id}", DeleteUserController},
}
