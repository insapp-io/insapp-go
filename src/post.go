package main

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Post struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `json:"name"`
	Association Association   `json:"association"`
	Description string        `json:"description"`
	Event       Event         `json:"event"`
	Date        time.Time     `json:"date"`
	Likes       Users         `json:"events_liked"`
	Comments    Comments      `json:"comments"`
	PhotoURL    string        `json:"photo_url"`
}

type Posts []Post

type Comment struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	User    User          `json:"user"`
	Content string        `json:"content"`
}

type Comments []Comment
