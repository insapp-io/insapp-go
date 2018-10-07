package main

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Comment defines how to model a Comment of a Post
type Comment struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	User    bson.ObjectId `json:"user"`
	Content string        `json:"content"`
	Date    time.Time     `json:"date"`
	Tags    Tags          `json:"tags"`
}

// Comments is an array of Comment
type Comments []Comment

type Tag struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	User string        `json:"user"`
	Name string        `json:"name"`
}

type Tags []Tag

// CommentPost will add the given comment object to the
// list of comments of the post linked to the given id
func CommentPost(id bson.ObjectId, comment Comment) Post {
	session := GetMongoSession()
	defer session.Close()
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
	session := GetMongoSession()
	defer session.Close()
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

func CommentEvent(id bson.ObjectId, comment Comment) Event {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("event")

	eventID := bson.M{"_id": id}
	change := bson.M{"$addToSet": bson.M{
		"comments": comment,
	}}

	db.Update(eventID, change)
	var event Event
	db.Find(bson.M{"_id": id}).One(&event)

	return event
}

func UncommentEvent(id bson.ObjectId, commentID bson.ObjectId) Event {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("event")

	DeleteNotificationsForComment(commentID)
	eventID := bson.M{"_id": id}
	change := bson.M{"$pull": bson.M{
		"comments": bson.M{"_id": commentID},
	}}

	db.Update(eventID, change)
	var event Event
	db.Find(bson.M{"_id": id}).One(&event)

	return event
}

func ReportComment(id bson.ObjectId, commentID bson.ObjectId, reporterId bson.ObjectId) {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("post")

	var post Post
	db.Find(bson.M{"_id": id}).One(&post)
	db = session.DB("insapp").C("user")
	var reporter User
	db.Find(bson.M{"_id": reporterId}).One(&reporter)
	for _, comment := range post.Comments {
		if comment.ID == commentID {
			var sender User
			db = session.DB("insapp").C("user")
			db.Find(bson.M{"_id": comment.User}).One(&sender)
			SendEmail("aeir@insa-rennes.fr", "Un commentaire a été reporté sur Insapp",
				"Ce commentaire a été reporté le "+time.Now().String()+
					"\n\nReporteur:\n"+reporter.ID.Hex()+"\n"+reporter.Username+
					"\n\nCommentaire:\n"+comment.ID.Hex()+"\n"+comment.Content+
					"\n\nPost:\n"+post.Title+
					"\n\nUser:\n"+sender.ID.Hex()+"\n"+sender.Username+"\n"+sender.Name)
		}
	}
}

func GetComment(postId bson.ObjectId, id bson.ObjectId) (Comment, error) {
	post := GetPost(postId)
	for _, comment := range post.Comments {
		if comment.ID == id {
			return comment, nil
		}
	}
	return Comment{}, errors.New("no comment found")
}

func GetCommentForEvent(eventId bson.ObjectId, id bson.ObjectId) (Comment, error) {
	event := GetEvent(eventId)
	for _, comment := range event.Comments {
		if comment.ID == id {
			return comment, nil
		}
	}
	return Comment{}, errors.New("no comment found")
}

func getCommentforUser(id bson.ObjectId, userId bson.ObjectId) []bson.ObjectId {
	post := GetPost(id)
	comments := post.Comments
	var results []bson.ObjectId
	for _, comment := range comments {
		if comment.User == userId {
			results = append(results, comment.ID)
		}
	}
	return results
}

func getCommentForUserOnEvent(id bson.ObjectId, userId bson.ObjectId) []bson.ObjectId {
	event := GetEvent(id)
	comments := event.Comments
	var results []bson.ObjectId
	for _, comment := range comments {
		if comment.User == userId {
			results = append(results, comment.ID)
		}
	}
	return results
}

func DeleteCommentsForUser(userId bson.ObjectId) {
	posts := GetLatestPosts(100)
	for _, post := range posts {
		comments := getCommentforUser(post.ID, userId)
		for _, commentId := range comments {
			UncommentPost(post.ID, commentId)
		}
	}
}

func DeleteCommentsForUserOnEvents(userId bson.ObjectId) {
	events := GetEvents()
	for _, event := range events {
		comments := getCommentForUserOnEvent(event.ID, userId)
		for _, commentId := range comments {
			UncommentEvent(event.ID, commentId)
		}
	}
}

func DeleteTagsForUserOnEvents(userId bson.ObjectId) {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("event")

	events := GetEvents()
	for _, event := range events {
		comments := event.Comments
		finalComments := Comments{}
		for _, comment := range comments {
			tags := comment.Tags
			finalTags := Tags{}
			for _, tag := range tags {
				if tag.User != userId.Hex() {
					finalTags = append(finalTags, tag)
				}
			}
			comment.Tags = finalTags
			finalComments = append(finalComments, comment)
		}
		db.Update(bson.M{"_id": event.ID}, bson.M{"$set": bson.M{"comments": finalComments}})
	}
}

func DeleteTagsForUser(userId bson.ObjectId) {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("post")

	var posts Posts
	db.Find(bson.M{}).All(&posts)
	for _, post := range posts {
		comments := post.Comments
		finalComments := Comments{}
		for _, comment := range comments {
			tags := comment.Tags
			finalTags := Tags{}
			for _, tag := range tags {
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
