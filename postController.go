package insapp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// GetPostController will answer a JSON of the post
// linked to the given id in the URL
func GetPostController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]
	var res = GetPost(bson.ObjectIdHex(postID))
	json.NewEncoder(w).Encode(res)
}

// GetAllPostsController will answer a JSON of the
// N latest posts. Here N = 50.
func GetAllPostsController(w http.ResponseWriter, r *http.Request) {
	id, err := GetUserFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(bson.M{"error": "could not get user ID"})
		return
	}

	user := GetUser(id)
	os := GetNotificationUserForUser(id).Os
	posts := GetLatestPosts(10)
	filteredPosts := Posts{}
	if user.ID != "" {
		for _, post := range posts {
			if contains(strings.ToUpper(user.Promotion), post.Promotions) && (contains(os, post.Plateforms) || os == "") || len(post.Promotions) == 0 || len(post.Plateforms) == 0 {
				filteredPosts = append(filteredPosts, post)
			}
		}
	} else {
		filteredPosts = posts
	}

	pagination := r.URL.Query().Get("range")
	if len(pagination) > 0 {
		re := regexp.MustCompile("\\[([0-9]+),\\s*?([0-9]+)\\]")
		matches := re.FindStringSubmatch(pagination)
		if len(matches) == 3 {
			if start, err := strconv.Atoi(matches[1]); err == nil {
				if end, err := strconv.Atoi(matches[2]); err == nil {
					if end >= start && len(filteredPosts) > start {
						paginatedPosts := Posts{}
						for i := start; i <= end && i < len(filteredPosts); i++ {
							paginatedPosts = append(paginatedPosts, filteredPosts[i])
						}
						json.NewEncoder(w).Encode(paginatedPosts)
						return
					}
				}
			}
		}
	}

	json.NewEncoder(w).Encode(filteredPosts)
}

// GetPostsForAssociationController will answer a JSON of the post owned by
// the given association
func GetPostsForAssociationController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	associationID := vars["id"]

	userID, err := GetUserFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(bson.M{"error": "could not get user ID"})
		return
	}

	user := GetUser(userID)
	os := GetNotificationUserForUser(userID).Os
	posts := GetPostsForAssociation(bson.ObjectIdHex(associationID))

	filteredPosts := Posts{}
	if user.ID != "" {
		for _, post := range posts {
			if contains(strings.ToUpper(user.Promotion), post.Promotions) && (contains(os, post.Plateforms) || os == "") || len(post.Promotions) == 0 || len(post.Plateforms) == 0 {
				filteredPosts = append(filteredPosts, post)
			}
		}
	} else {
		filteredPosts = posts
	}

	_ = json.NewEncoder(w).Encode(filteredPosts)
}

// AddPostController will answer a JSON of the
// brand new created post (from the JSON Body)
// Should be protected
func AddPostController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var post Post
	_ = decoder.Decode(&post)
	post.Date = time.Now()

	res := AddPost(post)
	association := GetAssociation(post.Association)
	_ = json.NewEncoder(w).Encode(res)
	go TriggerNotificationForPost(post, association.ID, res.ID, "@"+strings.ToLower(association.Name)+" a posté une news 📰")
}

// UpdatePostController will answer the JSON of the
// modified post (from the JSON Body)
// Should be protected
func UpdatePostController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var post Post
	_ = decoder.Decode(&post)
	vars := mux.Vars(r)
	postID := vars["id"]

	res := UpdatePost(bson.ObjectIdHex(postID), post)
	_ = json.NewEncoder(w).Encode(res)
}

// DeletePostController will answer a JSON of an
// empty post if the deletion has succeed
// Should be protected
func DeletePostController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post := GetPost(bson.ObjectIdHex(vars["id"]))

	res := DeletePost(post)
	_ = json.NewEncoder(w).Encode(res)
}

// LikePostController will answer a JSON of the
// post and the user that liked the post
func LikePostController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]
	userID := vars["userID"]

	post, user := LikePostWithUser(bson.ObjectIdHex(postID), bson.ObjectIdHex(userID))

	_ = json.NewEncoder(w).Encode(bson.M{"post": post, "user": user})
}

// DislikePostController will answer a JSON of the
// post and the user that disliked the post
func DislikePostController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]
	userID := vars["userID"]

	post, user := DislikePostWithUser(bson.ObjectIdHex(postID), bson.ObjectIdHex(userID))

	_ = json.NewEncoder(w).Encode(bson.M{"post": post, "user": user})
}

// CommentPostController will answer a JSON of the post
func CommentPostController(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "unable to read the request body"})
	}
	var comment Comment
	if err := json.Unmarshal([]byte(string(body)), &comment); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "wrong format"})
		return
	}

	comment.ID = bson.NewObjectId()
	comment.Date = time.Now()

	vars := mux.Vars(r)
	postID := vars["id"]

	post := CommentPost(bson.ObjectIdHex(postID), comment)
	association := GetAssociation(post.Association)
	user := GetUser(comment.User)

	_ = json.NewEncoder(w).Encode(post)

	if !post.NoNotification {
		_ = SendAssociationEmailForCommentOnPost(association.Email, post, comment, user)
	}

	for _, tag := range comment.Tags {
		go TriggerNotificationForUserFromPost(comment.User, bson.ObjectIdHex(tag.User), post.ID, "@"+GetUser(comment.User).Username+" t'a taggé sur '"+post.Title+"'", comment, "tag")
	}
}

// UncommentPostController will answer a JSON of the post
// Should be protected
func UncommentPostController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]
	commentID := vars["commentID"]

	res := UncommentPost(bson.ObjectIdHex(postID), bson.ObjectIdHex(commentID))

	_ = json.NewEncoder(w).Encode(res)
}

func ReportCommentController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]
	commentID := vars["commentID"]

	userID, err := GetUserFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(bson.M{"error": "could not get user ID"})
		return
	}

	ReportComment(bson.ObjectIdHex(postID), bson.ObjectIdHex(commentID), userID)
	json.NewEncoder(w).Encode(bson.M{})
}
