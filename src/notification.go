package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// NotificationUser defines how to model a NotificationUser
type NotificationUser struct {
	ID          bson.ObjectId   `bson:"_id,omitempty"`
	UserId      bson.ObjectId   `json:"userid"`
	Token       string          `json:"token"`
	Os          string          `json:"os"`
}


func CreateOrUpdateNotificationUser(user NotificationUser){
  session, _ := mgo.Dial("127.0.0.1")
  defer session.Close()
  session.SetMode(mgo.Monotonic, true)
  db := session.DB("insapp").C("notification")
  res, _ := db.Find(bson.M{"userid": user.UserId}).Count()
  if res > 0 {
  	db.Update(bson.M{"userid": user.UserId}, bson.M{"$set": bson.M{ "token": user.Token, "os": user.Os }})
  }else{
    db.Insert(user)
  }
}


func DeleteNotificationForUser(id bson.ObjectId){
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("notification")
	db.Remove(bson.M{"userid": id})
}
