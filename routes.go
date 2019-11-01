package insapp

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

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
func NewRouter(config Config) *mux.Router {
	err := InitJWT(config)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter().StrictSlash(true)

	for _, route := range publicRoutes {
		router.
			HandleFunc(route.Pattern, route.HandlerFunc).
			Methods(route.Method)
	}

	for _, route := range userRoutes {
		chain := chain(route.HandlerFunc, AuthMiddleware)
		router.
			HandleFunc(route.Pattern, chain).
			Methods(route.Method)
	}

	for _, route := range associationRoutes {
		chain := chain(route.HandlerFunc, AuthMiddleware)
		router.
			HandleFunc(route.Pattern, chain).
			Methods(route.Method)
	}

	for _, route := range superRoutes {
		chain := chain(route.HandlerFunc, AuthMiddleware)
		router.
			HandleFunc(route.Pattern, chain).
			Methods(route.Method)
	}

	return router
}

func chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

var publicRoutes = Routes{
	Route{"Index", "GET", "/", Index},
	Route{"HowToPost", "GET", "/how-to-post", HowToPost},
	Route{"Credit", "GET", "/credit", Credit},
	Route{"Legal", "GET", "/legal", Legal},
	/*
		Route{"LogAssociation", "POST", "/login/association", LogAssociationController},
		Route{"LogUser", "POST", "/login/user", LogUserController},
	*/
	Route{"SignUser", "POST", "/signin/user/{ticket}", LogInUserController},
}

var userRoutes = Routes{
	// Associations
	Route{"GetAssociation", "GET", "/associations", GetAllAssociationsController},
	Route{"GetAssociation", "GET", "/associations/{id}", GetAssociationController},
	Route{"GetPostsForAssociation", "GET", "/associations/{id}/posts", GetPostsForAssociationController},
	Route{"GetEventsForAssociation", "GET", "/associations/{id}/events", GetEventsForAssociationController},

	// Events
	Route{"GetFutureEvents", "GET", "/events", GetFutureEventsController},
	Route{"GetEvent", "GET", "/events/{id}", GetEventController},
	Route{"AddAttendee", "POST", "/events/{id}/attend/{userID}/status/{status}", ChangeAttendeeStatusController},
	Route{"RemoveAttendee", "DELETE", "/events/{id}/attend/{userID}", RemoveAttendeeController},
	Route{"CommentEvent", "POST", "/events/{id}/comment", CommentEventController},
	Route{"UncommentEvent", "DELETE", "/events/{id}/comment/{commentID}", UncommentEventController},

	// Posts
	Route{"GetPost", "GET", "/posts", GetAllPostsController},
	Route{"GetPost", "GET", "/posts/{id}", GetPostController},
	Route{"LikePost", "POST", "/posts/{id}/like/{userID}", LikePostController},
	Route{"DislikePost", "DELETE", "/posts/{id}/like/{userID}", DislikePostController},
	Route{"CommentPost", "POST", "/posts/{id}/comment", CommentPostController},
	Route{"UncommentPost", "DELETE", "/posts/{id}/comment/{commentID}", UncommentPostController},

	// Users
	Route{"GetUser", "GET", "/users/{id}", GetUserController},
	Route{"UpdateUser", "PUT", "/users/{id}", UpdateUserController},
	Route{"DeleteUser", "DELETE", "/users/{id}", DeleteUserController},

	// Notifications
	Route{"Notification", "POST", "/notifications", UpdateNotificationUserController},
	Route{"Notification", "GET", "/notifications/{userID}", GetNotificationController},
	Route{"Notification", "DELETE", "/notifications/{userID}/{id}", DeleteNotificationController},

	// Report
	Route{"ReportUser", "PUT", "/report/user/{id}", ReportUserController},
	Route{"ReportComment", "PUT", "/report/{id}/comment/{commentID}", ReportCommentController},

	// Search
	Route{"SearchUser", "POST", "/search/users", SearchUserController},
	Route{"SearchAssociation", "POST", "/search/associations", SearchAssociationController},
	Route{"SearchEvent", "POST", "/search/events", SearchEventController},
	Route{"SearchPost", "POST", "/search/posts", SearchPostController},
	Route{"SearchUniversal", "POST", "/search", SearchUniversalController},
}

var associationRoutes = Routes{
	// Associations
	Route{"UpdateAssociation", "PUT", "/associations/{id}", UpdateAssociationController},

	// Events
	Route{"AddEvent", "POST", "/events", AddEventController},
	Route{"UpdateEvent", "PUT", "/events/{id}", UpdateEventController},
	Route{"DeleteEvent", "DELETE", "/events/{id}", DeleteEventController},

	// Posts
	Route{"AddPost", "POST", "/posts", AddPostController},
	Route{"UpdatePost", "PUT", "/posts/{id}", UpdatePostController},
	Route{"DeletePost", "DELETE", "/posts/{id}", DeletePostController},

	// Image
	Route{"UploadNewImage", "POST", "/images", UploadNewImageController},
}

var superRoutes = Routes{
	// Users
	Route{"GetUsers", "GET", "/users", GetAllUserController},

	// Associations
	Route{"AddAssociation", "POST", "/associations", AddAssociationController},
	Route{"DeleteAssociation", "DELETE", "/associations/{id}", DeleteAssociationController},
	Route{"GetMyAssociations", "GET", "/associations/{id}/myassociations", GetMyAssociationController},
}
