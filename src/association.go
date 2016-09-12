package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Association defines the model of a Association
type Association struct {
	ID          bson.ObjectId   `bson:"_id,omitempty"`
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	Description string          `json:"description"`
	Events      []bson.ObjectId `json:"events"`
	Posts       []bson.ObjectId `json:"posts"`
	Cover    		string          `json:"profile"`
	Profile    	string          `json:"cover"`
	BgColor     string          `json:"bgColor"`
	FgColor     string          `json:"fgColor"`
}

// Associations is an array of Association
type Associations []Association

func AddAssociationUser(user AssociationUser) {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association_user")
	db.Insert(user)
}

// AddAssociation will add the given association to the database
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

// UpdateAssociation will update the given association link to the given ID,
// with the field of the given association, in the database
func UpdateAssociation(id bson.ObjectId, association Association) Association {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")
	assosID := bson.M{"_id": id}
	change := bson.M{"$set": bson.M{
		"name":        association.Name,
		"email":       association.Email,
		"description": association.Description,
		"profile":     association.Profile,
		"cover":     	 association.Cover,
		"bgcolor":     association.BgColor,
		"fgcolor":     association.FgColor,
	}}
	db.Update(assosID, change)
	var result Association
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}

// DeleteAssociation will delete the given association from the database
func DeleteAssociation(id bson.ObjectId) Association {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")
	db.RemoveId(id)
	var result Association
	db.FindId(id).One(result)
	return result
}

// GetAssociation will return an Association object from the given ID
func GetAssociation(id bson.ObjectId) Association {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")
	var result Association
	db.FindId(id).One(&result)
	return result
}

// GetAllAssociation will return an array of all the existing Association
func GetAllAssociation() Associations {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")
	var result Associations
	db.Find(bson.M{}).All(&result)
	return result
}

func GetMyAssociations(id bson.ObjectId) []bson.ObjectId {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association_user")
	var result []AssociationUser
	db.Find(bson.M{"owner": id}).All(&result)
	res := []bson.ObjectId{}
	for _, asso := range result {
		res = append(res, asso.Association)
	}
	return res
}

// AddEventToAssociation will add the given event ID to the given association
func AddEventToAssociation(id bson.ObjectId, event bson.ObjectId) Association {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")
	assosID := bson.M{"_id": id}
	change := bson.M{"$addToSet": bson.M{
		"events": event,
	}}
	db.Update(assosID, change)
	var result Association
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}

// RemoveEventFromAssociation will remove the given event ID from the given association
func RemoveEventFromAssociation(id bson.ObjectId, event bson.ObjectId) Association {
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

func AddPostToAssociation(id bson.ObjectId, post bson.ObjectId) Association {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")
	assosID := bson.M{"_id": id}
	change := bson.M{"$addToSet": bson.M{
		"posts": post,
	}}
	db.Update(assosID, change)
	var result Association
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}

func RemovePostFromAssociation(id bson.ObjectId, post bson.ObjectId) Association {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")
	assosID := bson.M{"_id": id}
	change := bson.M{"$pull": bson.M{
		"posts": post,
	}}
	db.Update(assosID, change)
	var result Association
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}

func SetImageAssociation(id bson.ObjectId, fileName string, isProfilePicture bool) Association {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")
	assosID := bson.M{"_id": id}
	if isProfilePicture {
		change := bson.M{"$set": bson.M{ "profile": fileName + ".png", }}
		db.Update(assosID, change)
	}else{
		change := bson.M{"$set": bson.M{ "cover": fileName + ".png", }}
		db.Update(assosID, change)
	}
	var result Association
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}
