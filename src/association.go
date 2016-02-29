package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Association struct {
	ID          bson.ObjectId   `bson:"_id,omitempty"`
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	Description string          `json:"description"`
	Events      []bson.ObjectId `json:"events"`
	PhotoURL    string          `json:"photoURL"`
	BgColor     string          `json:"bgColor"`
	FgColor     string          `json:"fgColor"`
}

type Associations []Association

func AddAssociation(association Association) Association {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")
	db.Insert(association)
	var result Association
	db.Find(bson.M{"name": association.Name}).One(&result)
	return result
}

func UpdateAssociation(id bson.ObjectId, association Association) Association {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")
	assosID := bson.M{"_id": id}
	change := bson.M{"$set": bson.M{
		"name":        association.Name,
		"Email":       association.Email,
		"Description": association.Description,
		"PohotURL":    association.PhotoURL,
		"BgColor":     association.BgColor,
		"FgColor":     association.FgColor,
	}}
	db.Update(assosID, change)
	var result Association
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}

func RemoveAssociation(id bson.ObjectId) Association {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")
	db.RemoveId(id)
	var result Association
	db.FindId(id).One(result)
	return result
}

func GetAssociation(id bson.ObjectId) Association {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")
	var result Association
	db.FindId(id).One(&result)
	return result
}

func GetAllAssociation() Associations {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")
	var result Associations
	db.Find(bson.M{}).All(&result)
	return result
}

func AddEventToAssociation(id bson.ObjectId, event bson.ObjectId) Association {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")
	assosID := bson.M{"_id": id}
	change := bson.M{"$push": bson.M{
		"events": event,
	}}
	db.Update(assosID, change)
	var result Association
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}

func RemoveEventToAssociation(id bson.ObjectId, event bson.ObjectId) Association {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")
	assosID := bson.M{"_id": id}
	change := bson.M{"$pull": bson.M{
		"events": event,
	}}
	db.Update(assosID, change)
	var result Association
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}
