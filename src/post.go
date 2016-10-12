package main

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Post defines how to model a Post
type Post struct {
	ID          bson.ObjectId   `bson:"_id,omitempty"`
	Title       string          `json:"title"`
	Association bson.ObjectId   `json:"association"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date"`
	Likes       []bson.ObjectId `json:"likes"`
	Comments    Comments        `json:"comments"`
	Image    		string          `json:"image"`
	ImageSize		bson.M					`json:"imageSize"`
}

// Posts is an array of Post
type Posts []Post

// Comment defines how to model a Comment of a Post
type Comment struct {
	ID      bson.ObjectId 		`bson:"_id,omitempty"`
	User    bson.ObjectId 		`json:"user"`
	Content string        		`json:"content"`
	Date    time.Time     		`json:"date"`
	Tags    Tags							`json:"tags"`
}

// Comments is an array of Comment
type Comments []Comment

type Tag struct {
	ID      bson.ObjectId 		`bson:"_id,omitempty"`
	User    string 						`json:"user"`
	Name 		string        		`json:"name"`
}

type Tags []Tag

// AddPost will add the given post to the database
func AddPost(post Post) Post {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("post")
	db.Insert(post)
	var result Post
	db.Find(bson.M{"title": post.Title, "date": post.Date}).One(&result)
	AddPostToAssociation(result.Association, result.ID)
	return result
}

// UpdatePost will update the post linked to the given ID,
// with the field of the given post, in the database
func UpdatePost(id bson.ObjectId, post Post) Post {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("post")
	postID := bson.M{"_id": id}
	change := bson.M{"$set": bson.M{
		"title"				:	post.Title,
		"description"	:	post.Description,
		"image"				:	post.Image,
		"imageSize"		:	post.ImageSize,
	}}
	db.Update(postID, change)
	var result Post
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}

// DeletePost will delete the given post from the database
func DeletePost(post Post) Post {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("post")
	db.RemoveId(post.ID)
	var result Post
	db.FindId(post.ID).One(result)
	DeleteNotificationsForPost(post.ID)
	RemovePostFromAssociation(post.Association, post.ID)
	for _, userId := range post.Likes{
		DislikePost(userId, post.ID)
	}
	return result
}

// GetPost will return an Post object from the given ID
func GetPost(id bson.ObjectId) Post {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("post")
	var result Post
	db.FindId(id).One(&result)
	return result
}

// GetLastestPosts will return an array of the last N Posts
func GetLastestPosts(number int) Posts {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("post")
	var result Posts
	db.Find(bson.M{}).Sort("-date").Limit(number).All(&result)
	return result
}

// LikePostWithUser will add the user to the list of
// user that liked the post (cf. Likes field)
func LikePostWithUser(id bson.ObjectId, userID bson.ObjectId) (Post, User) {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("post")
	postID := bson.M{"_id": id}
	change := bson.M{"$addToSet": bson.M{
		"likes": userID,
	}}
	db.Update(postID, change)
	var post Post
	db.Find(bson.M{"_id": id}).One(&post)
	user := LikePost(userID, post.ID)
	return post, user
}

// DislikePostWithUser will remove the user to the list of
// users that liked the post (cf. Likes field)
func DislikePostWithUser(id bson.ObjectId, userID bson.ObjectId) (Post, User) {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("post")
	postID := bson.M{"_id": id}
	change := bson.M{"$pull": bson.M{
		"likes": userID,
	}}
	db.Update(postID, change)
	var post Post
	db.Find(bson.M{"_id": id}).One(&post)
	user := DislikePost(userID, post.ID)
	return post, user
}

// CommentPost will add the given comment object to the
// list of comments of the post linked to the given id
func CommentPost(id bson.ObjectId, comment Comment) Post {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("post")
	postID := bson.M{"_id": id}
	change := bson.M{"$addToSet": bson.M{
		"comments": comment,
	}}
	db.Update(postID, change)
	var post Post
	db.Find(bson.M{"_id": id}).One(&post)
	return post
}

// UncommentPost will remove the given comment object from the
// list of comments of the post linked to the given id
func UncommentPost(id bson.ObjectId, commentID bson.ObjectId) Post {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("post")
	DeleteNotificationsForComment(commentID)
	postID := bson.M{"_id": id}
	change := bson.M{"$pull": bson.M{
		"comments": bson.M{"_id": commentID},
	}}
	db.Update(postID, change)
	var post Post
	db.Find(bson.M{"_id": id}).One(&post)
	return post
}

func ReportComment(id bson.ObjectId, commentID bson.ObjectId) {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("post")
	var post Post
	db.Find(bson.M{"_id": id}).One(&post)
	for _, comment := range post.Comments {
		if comment.ID == commentID {
			var sender User
			db = session.DB("insapp").C("user")
			db.Find(bson.M{"_id": id}).One(&sender)
			SendEmail("aeir@insa-rennes.fr", "Un commentaire a été reporté sur Insapp",
				"Ce commentaire a été reporté le " + time.Now().String() +
				"\n\nCommentaire:\n" + comment.ID.Hex() + "\n" + comment.Content +
				"\n\nPost:\n" + post.Title +
				"\n\nUser:\n" + sender.ID.Hex() + "\n" + sender.Username + "\n" + sender.Name)
		}
	}



}

func getCommentforUser(id bson.ObjectId, userId bson.ObjectId) []bson.ObjectId {
	post := GetPost(id)
	comments := post.Comments
	var results []bson.ObjectId
	for _, comment := range comments{
		if comment.User == userId {
			results = append(results, comment.ID)
		}
	}
	return results
}

func DeleteCommentsForUser(userId bson.ObjectId) {
	posts := GetLastestPosts(100)
	for _, post := range posts {
		comments := getCommentforUser(post.ID, userId)
		for _, commentId := range comments {
				UncommentPost(post.ID, commentId)
		}
	}
}

func DeleteTagsForUser(userId bson.ObjectId) {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("post")
	var posts Posts
	db.Find(bson.M{}).All(&posts)
	for _, post := range(posts){
		comments := post.Comments
		finalComments := Comments{}
		for _, comment := range(comments){
			tags := comment.Tags
			finalTags := Tags{}
			for _, tag := range(tags){
				if tag.User != userId.Hex() {
					finalTags = append(finalTags, tag)
				}
			}
			comment.Tags = finalTags
			finalComments = append(finalComments, comment)
		}
		db.Update(bson.M{"_id": post.ID}, bson.M{"$set": bson.M{"comments": finalComments}})
	}
}
