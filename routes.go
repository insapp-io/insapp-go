package insapp

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Middleware is the type wrapping http handlers.
type Middleware func(http.HandlerFunc, string) http.HandlerFunc

// Route type is used to define a route of the API
type Route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes type is an array of Route
type Routes []Route

// NewRouter is the constructor of the Router
// It will create every routes from the routes variable just above
func NewRouter() *mux.Router {
	err := InitJWT()
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
		router.
			HandleFunc(route.Pattern, AuthMiddleware(route.HandlerFunc, "user")).
			Methods(route.Method)
	}

	for _, route := range associationRoutes {
		router.
			HandleFunc(route.Pattern, AuthMiddleware(route.HandlerFunc, "association")).
			Methods(route.Method)
	}

	for _, route := range superRoutes {
		router.
			HandleFunc(route.Pattern, AuthMiddleware(route.HandlerFunc, "admin")).
			Methods(route.Method)
	}

	return router
}

var publicRoutes = Routes{
	Route{"GET", "/", Index},
	Route{"GET", "/how-to-post", HowToPost},
	Route{"GET", "/credit", Credit},
	Route{"GET", "/legal", Legal},

	// Login
	Route{"POST", "/login/user/{ticket}", LoginUserController},
	Route{"POST", "/login/association", LoginAssociationController},
}

var userRoutes = Routes{
	// Associations
	Route{"GET", "/associations", GetAllAssociationsController},
	Route{"GET", "/associations/{id}", GetAssociationController},
	Route{"GET", "/associations/{id}/posts", GetPostsForAssociationController},
	Route{"GET", "/associations/{id}/events", GetEventsForAssociationController},

	// Events
	Route{"GET", "/events", GetFutureEventsController},
	Route{"GET", "/events/{id}", GetEventController},

	Route{"POST", "/events/{id}/attend/{userID}/status/{status}", ChangeAttendeeStatusController},
	Route{"POST", "/events/{id}/comment", CommentEventController},

	Route{"DELETE", "/events/{id}/attend/{userID}", RemoveAttendeeController},
	Route{"DELETE", "/events/{id}/comment/{commentID}", UncommentEventController},

	// Posts
	Route{"GET", "/posts", GetAllPostsController},
	Route{"GET", "/posts/{id}", GetPostController},

	Route{"POST", "/posts/{id}/like/{userID}", LikePostController},
	Route{"POST", "/posts/{id}/comment", CommentPostController},

	Route{"DELETE", "/posts/{id}/like/{userID}", DislikePostController},
	Route{"DELETE", "/posts/{id}/comment/{commentID}", UncommentPostController},

	// Users
	Route{"GET", "/users/{id}", GetUserController},

	Route{"PUT", "/users/{id}", UpdateUserController},

	Route{"DELETE", "/users/{id}", DeleteUserController},

	// Notifications
	Route{"GET", "/notifications/{userID}", GetNotificationController},

	Route{"POST", "/notifications", UpdateNotificationUserController},

	Route{"DELETE", "/notifications/{userID}/{id}", DeleteNotificationController},

	// Report
	Route{"PUT", "/report/user/{id}", ReportUserController},
	Route{"PUT", "/report/{id}/comment/{commentID}", ReportCommentController},

	// Search
	Route{"POST", "/search/users", SearchUserController},
	Route{"POST", "/search/associations", SearchAssociationController},
	Route{"POST", "/search/events", SearchEventController},
	Route{"POST", "/search/posts", SearchPostController},
	Route{"POST", "/search", SearchUniversalController},

	// Logout
	Route{"POST", "/logout/user", LogoutUserController},
}

var associationRoutes = Routes{
	Route{"GET", "/association", GetAssociationUserController},

	// Associations
	Route{"PUT", "/associations/{id}", UpdateAssociationController},

	// Events
	Route{"POST", "/events", AddEventController},

	Route{"PUT", "/events/{id}", UpdateEventController},

	Route{"DELETE", "/events/{id}", DeleteEventController},

	// Posts
	Route{"POST", "/posts", AddPostController},

	Route{"PUT", "/posts/{id}", UpdatePostController},

	Route{"DELETE", "/posts/{id}", DeletePostController},

	// Image
	Route{"POST", "/images", UploadNewImageController},

	// Logout
	Route{"POST", "/logout/association", LogoutAssociationController},
}

var superRoutes = Routes{
	// Users
	Route{"GET", "/users", GetAllUserController},

	// Associations
	Route{"GET", "/associations/{id}/myassociations", GetMyAssociationController},

	Route{"POST", "/associations", AddAssociationController},

	Route{"DELETE", "/associations/{id}", DeleteAssociationController},
}
