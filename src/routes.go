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

// NewRouter is the constructor of the Router
// It will create every routes from the routes variable just above
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range publicRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	for _, route := range associationRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(tokenAuthAssociationUser.HandleFunc(route.HandlerFunc))
	}

	for _, route := range superRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(tokenAuthSuperUser.HandleFunc(route.HandlerFunc))
	}

	for _, route := range userRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(tokenAuthUser.HandleFunc(route.HandlerFunc))
	}

	return router
}

var tokenAuthAssociationUser = tauth.NewTokenAuth(nil, nil, memStoreAssociationUser, nil)
var tokenAuthSuperUser = tauth.NewTokenAuth(nil, nil, memStoreSuperUser, nil)
var tokenAuthUser = tauth.NewTokenAuth(nil, nil, memStoreUser, nil)

var memStoreAssociationUser = memstore.New("associationUser")
var memStoreSuperUser = memstore.New("superUser")
var memStoreUser = memstore.New("user")

var publicRoutes = Routes{
	Route{"Index", "GET", "/", Index},
	Route{"HowToPost", "GET", "/how-to-post", HowToPost},
	Route{"Credit", "GET", "/credit", Credit},
	Route{"Legal", "GET", "/legal", Legal},
	Route{"LogAssociation", "POST", "/login/association", LogAssociationController},
	Route{"LogUser", "POST", "/login/user", LogUserController},
	Route{"SignUser", "POST", "/signin/user/{ticket}", SignInUserController},
}

var superRoutes = Routes{
	Route{"GetUsers", "GET", "/users", GetAllUserController},
	Route{"AddAssociation", "POST", "/associations", AddAssociationController},
	Route{"DeleteAssociation", "DELETE", "/associations/{id}", DeleteAssociationController},
	Route{"GetMyAssociations", "GET", "/associations/{id}/myassociations", GetMyAssociationController},
}

var associationRoutes = Routes{
	//ASSOCIATIONS
	Route{"UpdateAssociation", "PUT", "/associations/{id}", UpdateAssociationController},

	//EVENTS
	Route{"AddEvent", "POST", "/events", AddEventController},
	Route{"UpdateEvent", "PUT", "/events/{id}", UpdateEventController},
	Route{"DeleteEvent", "DELETE", "/events/{id}", DeleteEventController},

	//POSTS
	Route{"AddPost", "POST", "/posts", AddPostController},
	Route{"UpdatePost", "PUT", "/posts/{id}", UpdatePostController},
	Route{"DeletePost", "DELETE", "/posts/{id}", DeletePostController},

	//IMAGE
	Route{"UploadNewImage", "POST", "/images", UploadNewImageController},
	Route{"UploadImage", "POST", "/images/{name}", UploadImageController},
}

var userRoutes = Routes{
	//ASSOCIATIONS
	Route{"GetAssociation", "GET", "/associations", GetAllAssociationsController},
	Route{"GetAssociation", "GET", "/associations/{id}", GetAssociationController},
	Route{"GetPostsForAssociation", "GET", "/associations/{id}/posts", GetPostsForAssociationController},
	Route{"GetEventsForAssociation", "GET", "/associations/{id}/events", GetEventsForAssociationController},

	//EVENTS
	Route{"GetFutureEvents", "GET", "/events", GetFutureEventsController},
	Route{"GetEvent", "GET", "/events/{id}", GetEventController},
	Route{"AddParticipant", "POST", "/events/{id}/attend/{userID}/status/{status}", ChangeAttendeeStatusController},
	Route{"RemoveParticipant", "DELETE", "/events/{id}/attend/{userID}", RemoveParticipantController},
	Route{"CommentEvent", "POST", "/events/{id}/comment", CommentEventController},
	Route{"UncommentEvent", "DELETE", "/events/{id}/comment/{commentID}", UncommentEventController},

	//POSTS
	Route{"GetPost", "GET", "/posts", GetAllPostsController},
	Route{"GetPost", "GET", "/posts/{id}", GetPostController},
	Route{"LikePost", "POST", "/posts/{id}/like/{userID}", LikePostController},
	Route{"DislikePost", "DELETE", "/posts/{id}/like/{userID}", DislikePostController},
	Route{"CommentPost", "POST", "/posts/{id}/comment", CommentPostController},
	Route{"UncommentPost", "DELETE", "/posts/{id}/comment/{commentID}", UncommentPostController},

	//USERS
	Route{"GetUser", "GET", "/users/{id}", GetUserController},
	Route{"UpdateUser", "PUT", "/users/{id}", UpdateUserController},
	Route{"DeleteUser", "DELETE", "/users/{id}", DeleteUserController},

	//NOTIFICATIONS
	Route{"Notification", "POST", "/notifications", UpdateNotificationUserController},
	Route{"Notification", "GET", "/notifications/{userID}", GetNotificationController},
	Route{"Notification", "DELETE", "/notifications/{userID}/{id}", DeleteNotificationController},

	//REPORTING
	Route{"ReportUser", "PUT", "/report/user/{id}", ReportUserController},
	Route{"ReportComment", "PUT", "/report/{id}/comment/{commentID}", ReportCommentController},

	//SEARCHING
	Route{"SearchUser", "POST", "/search/users", SearchUserController},
	Route{"SearchAssociation", "POST", "/search/associations", SearchAssociationController},
	Route{"SearchEvent", "POST", "/search/events", SearchEventController},
	Route{"SearchPost", "POST", "/search/posts", SearchPostController},
	Route{"SearchUniversal", "POST", "/search", SearchUniversalController},
}
