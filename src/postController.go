package main

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
)

// GetPostController will answer a JSON of the post
// linked to the given id in the URL
func GetPostController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]
	var res = GetPost(bson.ObjectIdHex(postID))
	json.NewEncoder(w).Encode(res)
}

// GetLastestPostsController will answer a JSON of the
// N lastest post. Here N = 50.
func GetLastestPostsController(w http.ResponseWriter, r *http.Request) {
	var res = GetLastestPosts(50)
	json.NewEncoder(w).Encode(res)
}

// AddPostController will answer a JSON of the
// brand new created post (from the JSON Body)
func AddPostController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var post Post
	decoder.Decode(&post)
	res := AddPost(post)
	json.NewEncoder(w).Encode(res)
}

// UpdatePostController will answer the JSON of the
// modified post (from the JSON Body)
func UpdatePostController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var post Post
	decoder.Decode(&post)
	vars := mux.Vars(r)
	postID := vars["id"]
	res := UpdatePost(bson.ObjectIdHex(postID), post)
	json.NewEncoder(w).Encode(res)
}

// DeletePostController will answer a JSON of an
// empty post if the deletation has succeed
func DeletePostController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	res := DeletePost(bson.ObjectIdHex(vars["id"]))
	json.NewEncoder(w).Encode(res)
}

// LikePostController will answer a JSON of the
// post and the user that liked the post
func LikePostController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]
	userID := vars["userID"]
	res := LikePost(bson.ObjectIdHex(postID), bson.ObjectIdHex(userID))
	json.NewEncoder(w).Encode(res)
}

// DislikePostController will answer a JSON of the
// post and the user that disliked the post
func DislikePostController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]
	userID := vars["userID"]
	res := DislikePost(bson.ObjectIdHex(postID), bson.ObjectIdHex(userID))
	json.NewEncoder(w).Encode(res)
}

// CommentPostController will answer a JSON of the post
func CommentPostController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var comment Comment
	decoder.Decode(&comment)
	vars := mux.Vars(r)
	postID := vars["id"]
	res := CommentPost(bson.ObjectIdHex(postID), comment)
	json.NewEncoder(w).Encode(res)
}

// UncommentPostController will answer a JSON of the post
func UncommentPostController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var comment Comment
	decoder.Decode(&comment)
	vars := mux.Vars(r)
	postID := vars["id"]
	res := UncommentPost(bson.ObjectIdHex(postID), comment)
	json.NewEncoder(w).Encode(res)
}
