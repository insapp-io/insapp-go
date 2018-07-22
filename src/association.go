package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Association defines the model of a Association
type Association struct {
	ID            bson.ObjectId   `bson:"_id,omitempty"`
	Name          string          `json:"name"`
	Email         string          `json:"email"`
	Description   string          `json:"description"`
	Events        []bson.ObjectId `json:"events"`
	Posts         []bson.ObjectId `json:"posts"`
	Palette       [][]int         `json:"palette"`
	SelectedColor int             `json:"selectedcolor"`
	Profile       string          `json:"profile"`
	Cover         string          `json:"cover"`
	BgColor       string          `json:"bgcolor"`
	FgColor       string          `json:"fgcolor"`
}

// Associations is an array of Association
type Associations []Association

func AddAssociationUser(user AssociationUser) {
	_, info, _ := Configuration()
	session, _ := mgo.DialWithInfo(info)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	db := session.DB("insapp").C("association_user")
	db.Insert(user)
}

// AddAssociation will add the given association to the database
func AddAssociation(association Association) Association {
	_, info, _ := Configuration()
	session, _ := mgo.DialWithInfo(info)
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
	_, info, _ := Configuration()
	session, _ := mgo.DialWithInfo(info)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")

	assosID := bson.M{"_id": id}
	change := bson.M{"$set": bson.M{
		"name":          association.Name,
		"email":         association.Email,
		"description":   association.Description,
		"profile":       association.Profile,
		"cover":         association.Cover,
		"palette":       association.Palette,
		"selectedcolor": association.SelectedColor,
		"bgcolor":       association.BgColor,
		"fgcolor":       association.FgColor,
	}}

	db.Update(assosID, change)
	var result Association
	db.Find(bson.M{"_id": id}).One(&result)

	return result
}

// DeleteAssociation will delete the given association from the database
func DeleteAssociation(id bson.ObjectId) Association {
	_, info, _ := Configuration()
	session, _ := mgo.DialWithInfo(info)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")

	association := GetAssociation(id)
	for _, eventId := range association.Events {
		DeleteEvent(GetEvent(eventId))
	}
	for _, postId := range association.Posts {
		DeletePost(GetPost(postId))
	}

	db.RemoveId(id)
	var result Association
	db.FindId(id).One(result)

	return result
}

// GetAssociation will return an Association object from the given ID
func GetAssociation(id bson.ObjectId) Association {
	_, info, _ := Configuration()
	session, _ := mgo.DialWithInfo(info)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")

	var result Association
	db.FindId(id).One(&result)

	return result
}

// GetAllAssociation will return an array of all the existing Association
func GetAllAssociation() Associations {
	_, info, _ := Configuration()
	session, _ := mgo.DialWithInfo(info)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")

	var result Associations
	db.Find(bson.M{}).All(&result)

	return result
}

func GetMyAssociations(id bson.ObjectId) []bson.ObjectId {
	_, info, _ := Configuration()
	session, _ := mgo.DialWithInfo(info)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association_user")

	var result []AssociationUser
	db.Find(bson.M{"owner": id}).All(&result)
	var res []bson.ObjectId
	for _, association := range result {
		res = append(res, association.Association)
	}

	return res
}

func SearchAssociation(name string) Associations {
	_, info, _ := Configuration()
	session, _ := mgo.DialWithInfo(info)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association")

	var result Associations
	db.Find(bson.M{"$or": []interface{}{
		bson.M{"name": bson.M{"$regex": bson.RegEx{`^.*` + name + `.*`, "i"}}}, bson.M{"description": bson.M{"$regex": bson.RegEx{`^.*` + name + `.*`, "i"}}}}}).All(&result)

	return result
}

// AddEventToAssociation will add the given event ID to the given association
func AddEventToAssociation(id bson.ObjectId, event bson.ObjectId) Association {
	_, info, _ := Configuration()
	session, _ := mgo.DialWithInfo(info)
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
	_, info, _ := Configuration()
	session, _ := mgo.DialWithInfo(info)
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
	_, info, _ := Configuration()
	session, _ := mgo.DialWithInfo(info)
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
	_, info, _ := Configuration()
	session, _ := mgo.DialWithInfo(info)
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

func GetAssociationUser(id bson.ObjectId) AssociationUser {
	_, info, _ := Configuration()
	session, _ := mgo.DialWithInfo(info)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association_user")

	var result AssociationUser
	db.Find(bson.M{"association": id}).One(&result)

	return result
}
