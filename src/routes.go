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
	Route{"Credit", "GET", "/credit", Credit},
	Route{"Legal", "GET", "/legal", Legal},
	Route{"LogAssociation", "POST", "/login/association", LogAssociationController},
	Route{"LogUser", "POST", "/login/user", LogUserController},
	Route{"SignUser", "POST", "/signin/user/{ticket}", SignInUserController},
}

var superRoutes = Routes{
	Route{"GetUsers", "GET", "/user", GetAllUserController},
	Route{"AddAssociation", "POST", "/association", AddAssociationController},
	Route{"DeleteAssociation", "DELETE", "/association/{id}", DeleteAssociationController},
	Route{"GetMyAssociations", "GET", "/association/{id}/myassociations", GetMyAssociationController},
}

var associationRoutes = Routes{
	//ASSOCIATIONS
	Route{"UpdateAssociation", "PUT", "/association/{id}", UpdateAssociationController},

	//EVENTS
	Route{"AddEvent", "POST", "/event", AddEventController},
	Route{"UpdateEvent", "PUT", "/event/{id}", UpdateEventController},
	Route{"DeleteEvent", "DELETE", "/event/{id}", DeleteEventController},

	//POSTS
	Route{"AddPost", "POST", "/post", AddPostController},
	Route{"UpdatePost", "PUT", "/post/{id}", UpdatePostController},
	Route{"DeletePost", "DELETE", "/post/{id}", DeletePostController},

	//IMAGE
	Route{"UploadNewImage", "POST", "/image", UploadNewImageController},
	Route{"UploadImage", "POST", "/image/{name}", UploadImageController},
}

var userRoutes = Routes{
	//ASSOCIATIONS
	Route{"GetAssociation", "GET", "/association", GetAllAssociationsController},
	Route{"GetAssociation", "GET", "/association/{id}", GetAssociationController},

	//EVENTS
	Route{"GetFutureEvents", "GET", "/event", GetFutureEventsController},
	Route{"GetEvent", "GET", "/event/{id}", GetEventController},
	Route{"AddParticipant", "POST", "/event/{id}/participant/{userID}/status/{status}", ChangeAttendeeStatusController},
	Route{"RemoveParticipant", "DELETE", "/event/{id}/participant/{userID}", RemoveParticipantController},
	Route{"CommentEvent", "POST", "/event/{id}/comment", CommentEventController},
	Route{"UncommentEvent", "DELETE", "/event/{id}/comment/{commentID}", UncommentEventController},

	//POSTS
	Route{"GetPost", "GET", "/post", GetAllPostsController},
	Route{"GetPost", "GET", "/post/{id}", GetPostController},
	Route{"LikePost", "POST", "/post/{id}/like/{userID}", LikePostController},
	Route{"DislikePost", "DELETE", "/post/{id}/like/{userID}", DislikePostController},
	Route{"CommentPost", "POST", "/post/{id}/comment", CommentPostController},
	Route{"UncommentPost", "DELETE", "/post/{id}/comment/{commentID}", UncommentPostController},
	Route{"ReportComment", "PUT", "/report/{id}/comment/{commentID}", ReportCommentController},

	//USERS
	Route{"GetUser", "GET", "/user/{id}", GetUserController},
	Route{"UpdateUser", "PUT", "/user/{id}", UpdateUserController},
	Route{"DeleteUser", "DELETE", "/user/{id}", DeleteUserController},
	Route{"ReportUser", "PUT", "/report/user/{id}", ReportUserController},

	//NOTIFICATIONS
	Route{"Notification", "POST", "/notification", UpdateNotificationUserController},
	Route{"Notification", "GET", "/notification/{userID}", GetNotificationController},
	Route{"Notification", "DELETE", "/notification/{userID}/{id}", DeleteNotificationController},

	//SEARCH
	Route{"SearchUser", "POST", "/search/users", SearchUserController},
	Route{"SearchAssociation", "POST", "/search/associations", SearchAssociationController},
	Route{"SearchEvent", "POST", "/search/events", SearchEventController},
	Route{"SearchPost", "POST", "/search/posts", SearchPostController},
	Route{"SearchUniversal", "POST", "/search", SearchUniversalController},
}
