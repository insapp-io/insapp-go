package main

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Event struct {
	ID           bson.ObjectId   `bson:"_id,omitempty"`
	Name         string          `json:"name"`
	Association  bson.ObjectId   `json:"association" bson:"association"`
	Description  string          `json:"description"`
	Participants []bson.ObjectId `json:"participants" bson:"participants,omitempty"`
	Status       string          `json:"status"`
	DateStart    time.Time       `json:"dateStart"`
	DateEnd      time.Time       `json:"dateEnd"`
	PhotoURL     string          `json:"photoURL"`
	BgColor      string          `json:"bgColor"`
	FgColor      string          `json:"fgColor"`
}

type Events []Event

func GetEvent(id bson.ObjectId) Event {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	var result Event
	db.FindId(id).One(&result)
	return result
}

func GetFutureEvent() Events {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	var result Events
	var now = time.Now()
	db.Find(bson.M{"datestart": bson.M{"$gt": now}}).All(&result)
	return result
}

func AddEvent(event Event) Event {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	db.Insert(event)
	var result Event
	db.Find(bson.M{"name": event.Name, "datestart": event.DateStart}).One(&result)
	AddEventToAssociation(result.Association, result.ID)
	return result
}

func UpdateEvent(id bson.ObjectId, event Event) Event {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	eventID := bson.M{"_id": id}
	change := bson.M{"$set": bson.M{
		"name":        event.Name,
		"description": event.Description,
		"status":      event.Status,
		"photoURL":    event.PhotoURL,
		"datestart":   event.DateStart,
		"dateend":     event.DateEnd,
		"bgcolor":     event.BgColor,
		"fgcolor":     event.FgColor,
	}}
	db.Update(eventID, change)
	var result Event
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}

func DeleteEvent(event Event) Event {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	db.Remove(event)
	RemoveEventFromAssociation(event.Association, event.ID)
	var result Event
	db.Find(event.ID).One(result)
	return result
}

func AddParticipant(id bson.ObjectId, userID bson.ObjectId) Event {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	eventID := bson.M{"_id": id}
	change := bson.M{"$push": bson.M{
		"participants": userID,
	}}
	db.Update(eventID, change)
	var result Event
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}

func RemoveParticipant(id bson.ObjectId, userID bson.ObjectId) Event {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	eventID := bson.M{"_id": id}
	change := bson.M{"$pull": bson.M{
		"participants": userID,
	}}
	db.Update(eventID, change)
	var result Event
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}
