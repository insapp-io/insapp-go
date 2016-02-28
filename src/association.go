package main

import "gopkg.in/mgo.v2/bson"

type Association struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `json:"name"`
	Email       string        `json:"email"`
	Description string        `json:"description"`
	Events      Events        `json:"events"`
	PhotoURL    string        `json:"photo_url"`
	BgColor     []uint8       `json:"bg_color"`
	FgColor     []uint8       `json:"fg_color"`
}

type Associations []Association
