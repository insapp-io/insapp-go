package main

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `json:"name"`
	Username    string        `json:"username"`
	Description string        `json:"description"`
	Email       string        `json:"email"`
	Promotion   string        `json:"promotion"`
	Events      Events        `json:"events_liked"`
	PostLiked   Posts         `json:"posts_liked"`
	PhotoURL    string        `json:"photo_url"`
}

type Users []User
