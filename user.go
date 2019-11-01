package insapp

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

var promotions = []string{"", "1STPI", "2STPI",
	"3CDTI", "3EII", "3GM", "3GCU", "3GMA", "3INFO", "3SGM", "3SRC",
	"4CDTI", "4EII", "4GM", "4GCU", "4GMA", "4INFO", "4SGM", "4SRC",
	"5CDTI", "5EII", "5GM", "5GCU", "5GMA", "5INFO", "5SGM", "5SRC",
	"STAFF"}

var genders = []string{"", "female", "male"}

// User defines how to model a User
type User struct {
	ID          bson.ObjectId   `bson:"_id,omitempty"`
	Name        string          `json:"name"`
	Username    string          `json:"username"`
	Description string          `json:"description"`
	Email       string          `json:"email"`
	EmailPublic bool            `json:"emailpublic"`
	Promotion   string          `json:"promotion"`
	Gender      string          `json:"gender"`
	Events      []bson.ObjectId `json:"events"`
	PostsLiked  []bson.ObjectId `json:"postsliked"`
}

// Users is an array of User
type Users []User

// AddUser will add the given user from JSON body to the database
func AddUser(user User) User {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("user")

	db.Insert(user)
	var result User
	db.Find(bson.M{"username": user.Username}).One(&result)

	return result
}

// UpdateUser will update the user link to the given ID,
// with the field of the given user, in the database
func UpdateUser(id bson.ObjectId, user User) User {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("user")

	promotion := ""
	for _, promo := range promotions {
		if user.Promotion == promo {
			promotion = promo
			break
		}
	}

	gender := ""
	for _, gen := range genders {
		if user.Gender == gen {
			gender = gen
			break
		}
	}

	userID := bson.M{"_id": id}
	change := bson.M{"$set": bson.M{
		"name":        user.Name,
		"description": user.Description,
		"email":       user.Email,
		"emailpublic": user.EmailPublic,
		"promotion":   promotion,
		"gender":      gender,
	}}
	db.Update(userID, change)

	var result User
	db.Find(bson.M{"_id": id}).One(&result)

	return result
}

// DeleteUser will delete the given user from the database
func DeleteUser(user User) User {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("user")

	//DeleteCredentialsForUser(user.ID)
	DeleteNotificationsForUser(user.ID)
	DeleteNotificationTokenForUser(user.ID)

	for _, eventID := range user.Events {
		RemoveAttendee(eventID, user.ID, "going")
		RemoveAttendee(eventID, user.ID, "notgoing")
		RemoveAttendee(eventID, user.ID, "maybe")
	}

	for _, postID := range user.PostsLiked {
		DislikePostWithUser(postID, user.ID)
	}

	DeleteTagsForUser(user.ID)
	DeleteTagsForUserOnEvents(user.ID)
	DeleteCommentsForUser(user.ID)
	DeleteCommentsForUserOnEvents(user.ID)
	db.RemoveId(user.ID)

	var result User
	db.FindId(user.ID).One(result)

	return result
}

// GetAllUser will return an User object from the given ID
func GetAllUser() Users {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("user")

	var result Users
	db.Find(bson.M{}).All(&result)

	return result
}

// GetUser will return an User object from the given ID
func GetUser(id bson.ObjectId) User {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("user")

	var result User
	db.FindId(id).One(&result)

	return result
}

// LikePost will add the postID to the list of liked post
// of the user linked to the given id
func LikePost(id bson.ObjectId, postID bson.ObjectId) User {
	session := GetMongoSession()
	defer session.Close()
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
	session := GetMongoSession()
	defer session.Close()
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
	session := GetMongoSession()
	defer session.Close()
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
	session := GetMongoSession()
	defer session.Close()
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

func SearchUser(name string) Users {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("user")

	var result Users
	db.Find(bson.M{"$or": []interface{}{
		bson.M{"username": bson.M{"$regex": bson.RegEx{`^.*` + name + `.*`, "i"}}}, bson.M{"name": bson.M{"$regex": bson.RegEx{`^.*` + name + `.*`, "i"}}}}}).All(&result)

	return result
}

func ReportUser(id bson.ObjectId, reporterID bson.ObjectId) {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("user")

	var user User
	db.Find(bson.M{"_id": id}).One(&user)
	var reporter User
	db.Find(bson.M{"_id": reporterID}).One(&reporter)
	SendEmail("aeir@insa-rennes.fr", "Un utilisateur a été reporté sur Insapp",
		"Cet utilisateur a été reporté le "+time.Now().String()+
			"\n\nReporteur:\n"+reporter.ID.Hex()+"\n"+reporter.Username+"\n"+reporter.Name+
			"\n\nSignaler:\n"+user.ID.Hex()+"\n"+user.Username+"\n"+user.Name+"\n"+user.Description)
}
