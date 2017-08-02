package main

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Event defines what an Event is
type Event struct {
	ID             bson.ObjectId   `bson:"_id,omitempty"`
	Name           string          `json:"name"`
	Association    bson.ObjectId   `json:"association" bson:"association"`
	Description    string          `json:"description"`
	Participants   []bson.ObjectId `json:"participants" bson:"participants,omitempty"`
	Maybe          []bson.ObjectId `json:"maybe" bson:"maybe,omitempty"`
	NotGoing       []bson.ObjectId `json:"notgoing" bson:"notgoing,omitempty"`
	Comments       Comments        `json:"comments"`
	Status         string          `json:"status"`
	Palette        [][]int         `json:"palette"`
	SelectedColor  int             `json:"selectedcolor"`
	DateStart      time.Time       `json:"dateStart"`
	DateEnd        time.Time       `json:"dateEnd"`
	Image          string          `json:"image"`
	Promotions     []string        `json:"promotions"`
	Plateforms     []string        `json:"plateforms"`
	BgColor        string          `json:"bgColor"`
	FgColor        string          `json:"fgColor"`
	NoNotification bool            `json:"nonotification"`
}

// Events is an array of Event
type Events []Event

// GetEvent returns an Event object from the given ID
func GetEvent(id bson.ObjectId) Event {
	conf, _ := Configuration()
	session, _ := mgo.Dial(conf.Database)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	var result Event
	db.FindId(id).One(&result)
	return result
}

func GetEvents() Events {
	conf, _ := Configuration()
	session, _ := mgo.Dial(conf.Database)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	var result Events
	db.Find(bson.M{}).All(&result)
	return result
}

// GetFutureEvents returns an array of Event objects
// that will happen after "NOW"
func GetFutureEvents() Events {
	conf, _ := Configuration()
	session, _ := mgo.Dial(conf.Database)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	var result Events
	var now = time.Now()
	db.Find(bson.M{"dateend": bson.M{"$gt": now}}).All(&result)
	return result
}

// AddEvent will add the Event event to the database
func AddEvent(event Event) Event {
	conf, _ := Configuration()
	session, _ := mgo.Dial(conf.Database)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	db.Insert(event)
	var result Event
	db.Find(bson.M{"name": event.Name, "datestart": event.DateStart}).One(&result)
	AddEventToAssociation(result.Association, result.ID)
	return result
}

// UpdateEvent will update the Event event in the database
func UpdateEvent(id bson.ObjectId, event Event) Event {
	conf, _ := Configuration()
	session, _ := mgo.Dial(conf.Database)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	eventID := bson.M{"_id": id}
	change := bson.M{"$set": bson.M{
		"name":           event.Name,
		"description":    event.Description,
		"status":         event.Status,
		"image":          event.Image,
		"palette":        event.Palette,
		"selectedcolor":  event.SelectedColor,
		"datestart":      event.DateStart,
		"dateend":        event.DateEnd,
		"plateforms":     event.Plateforms,
		"promotions":     event.Promotions,
		"bgcolor":        event.BgColor,
		"fgcolor":        event.FgColor,
		"nonotification": event.NoNotification,
	}}
	db.Update(eventID, change)
	var result Event
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}

// DeleteEvent will delete the given Event
func DeleteEvent(event Event) Event {
	conf, _ := Configuration()
	session, _ := mgo.Dial(conf.Database)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	db.Remove(event)
	DeleteNotificationsForEvent(event.ID)
	RemoveEventFromAssociation(event.Association, event.ID)
	for _, userId := range event.Participants {
		RemoveEventFromUser(userId, event.ID)
	}
	var result Event
	db.Find(event.ID).One(result)
	return result
}

// AddParticipant add the given userID to the given eventID as a participant
func AddParticipantToGoingList(id bson.ObjectId, userID bson.ObjectId) (Event, User) {
	RemoveParticipant(id, userID, "notgoing")
	RemoveParticipant(id, userID, "maybe")
	conf, _ := Configuration()
	session, _ := mgo.Dial(conf.Database)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	eventID := bson.M{"_id": id}
	change := bson.M{"$addToSet": bson.M{
		"participants": userID,
	}}
	db.Update(eventID, change)
	var event Event
	db.Find(bson.M{"_id": id}).One(&event)
	user := AddEventToUser(userID, event.ID)
	return event, user
}

func AddParticipantToMaybeList(id bson.ObjectId, userID bson.ObjectId) (Event, User) {
	RemoveParticipant(id, userID, "notgoing")
	RemoveParticipant(id, userID, "participants")
	conf, _ := Configuration()
	session, _ := mgo.Dial(conf.Database)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	eventID := bson.M{"_id": id}
	change := bson.M{"$addToSet": bson.M{
		"maybe": userID,
	}}
	db.Update(eventID, change)
	var event Event
	db.Find(bson.M{"_id": id}).One(&event)
	user := GetUser(userID)
	return event, user
}

func AddParticipantToNotGoingList(id bson.ObjectId, userID bson.ObjectId) (Event, User) {
	RemoveParticipant(id, userID, "maybe")
	RemoveParticipant(id, userID, "participants")
	conf, _ := Configuration()
	session, _ := mgo.Dial(conf.Database)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	eventID := bson.M{"_id": id}
	change := bson.M{"$addToSet": bson.M{
		"notgoing": userID,
	}}
	db.Update(eventID, change)
	var event Event
	db.Find(bson.M{"_id": id}).One(&event)
	user := GetUser(userID)
	return event, user
}

// RemoveParticipant remove the given userID from the given eventID as a participant
func RemoveParticipant(id bson.ObjectId, userID bson.ObjectId, list string) (Event, User) {
	conf, _ := Configuration()
	session, _ := mgo.Dial(conf.Database)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	eventID := bson.M{"_id": id}
	change := bson.M{"$pull": bson.M{
		list: userID,
	}}
	db.Update(eventID, change)
	var event Event
	db.Find(bson.M{"_id": id}).One(&event)
	user := RemoveEventFromUser(userID, event.ID)
	return event, user
}

func SearchEvent(name string) Events {
	conf, _ := Configuration()
	session, _ := mgo.Dial(conf.Database)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	var result Events
	db.Find(bson.M{"$or": []interface{}{
		bson.M{"name": bson.M{"$regex": bson.RegEx{`^.*` + name + `.*`, "i"}}}, bson.M{"description": bson.M{"$regex": bson.RegEx{`^.*` + name + `.*`, "i"}}}}}).All(&result)
	return result
}

//CommentEvent(eventId, comment)
//UncommentEvent(eventId, comment)
//GetCommentFromEvent(eventId, comment)
