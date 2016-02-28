package main

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID          bson.ObjectId   `bson:"_id,omitempty"`
	Name        string          `json:"name"`
	Username    string          `json:"username"`
	Description string          `json:"description"`
	Email       string          `json:"email"`
	Promotion   string          `json:"promotion"`
	Events      []bson.ObjectId `json:"events_liked"`
	PostLiked   []bson.ObjectId `json:"posts_liked"`
	PhotoURL    string          `json:"photo_url"`
}

type Users []User
