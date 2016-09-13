package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// User defines how to model a User
type User struct {
	ID          bson.ObjectId   `bson:"_id,omitempty"`
	Name        string          `json:"name"`
	Username    string          `json:"username"`
	Description string          `json:"description"`
	Email       string          `json:"email"`
	EmailPublic bool            `json:"emailpublic"`
	Promotion   string          `json:"promotion"`
	Events      []bson.ObjectId `json:"events"`
	PostsLiked  []bson.ObjectId `json:"postsliked"`
}

// Users is an array of User
type Users []User

// AddUser will add the given user from JSON body to the database
func AddUser(user User) User {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("user")
	db.Insert(user)
	var result User
	db.Find(bson.M{"username": user.Username}).One(&result)
	return result
}

// UpdateUser will update the user link to the given ID,
// with the field of the given user, in the database
func UpdateUser(id bson.ObjectId, user User) User {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("user")
	userID := bson.M{"_id": id}
	change := bson.M{"$set": bson.M{
		"name":        user.Name,
		"description": user.Description,
		"emailpublic": user.EmailPublic,
		"promotion":   user.Promotion,
	}}
	db.Update(userID, change)
	var result User
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}

// DeleteUser will delete the given user from the database
func DeleteUser(id bson.ObjectId) User {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("user")
	db.RemoveId(id)
	var result User
	db.FindId(id).One(result)
	return result
}

// GetUser will return an User object from the given ID
func GetUser(id bson.ObjectId) User {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("user")
	var result User
	db.FindId(id).One(&result)
	return result
}

// LikePost will add the postID to the list of liked post
// of the user linked to the given id
func LikePost(id bson.ObjectId, postID bson.ObjectId) User {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("user")
	userID := bson.M{"_id": id}
	change := bson.M{"$addToSet": bson.M{
		"postsliked": postID,
	}}
	db.Update(userID, change)
	var result User
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}

// DislikePost will remove the postID from the list of liked
// post of the user linked to the given id
func DislikePost(id bson.ObjectId, postID bson.ObjectId) User {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("user")
	userID := bson.M{"_id": id}
	change := bson.M{"$pull": bson.M{
		"postsliked": postID,
	}}
	db.Update(userID, change)
	var result User
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}

// AddEventToUser will add the eventID to the list
// of the user's event linked to the given id
func AddEventToUser(id bson.ObjectId, eventID bson.ObjectId) User {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("user")
	userID := bson.M{"_id": id}
	change := bson.M{"$addToSet": bson.M{
		"events": eventID,
	}}
	db.Update(userID, change)
	var result User
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}

// RemoveEventFromUser will remove the eventID from the list
// of the user's event linked to the given id
func RemoveEventFromUser(id bson.ObjectId, eventID bson.ObjectId) User {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("user")
	userID := bson.M{"_id": id}
	change := bson.M{"$pull": bson.M{
		"events": eventID,
	}}
	db.Update(userID, change)
	var result User
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}

func SetImageUser(id bson.ObjectId, fileName string) User {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("user")
	userID := bson.M{"_id": id}
	change := bson.M{
		"photoUrl": fileName,
	}
	db.Update(userID, change)
	var result User
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}
