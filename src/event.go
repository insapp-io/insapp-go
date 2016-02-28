package main

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Event struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Name         string        `json:"name"`
	Association  Association   `json:"association"`
	Description  string        `json:"description"`
	Participants Users         `json:"participants"`
	DateStart    time.Time     `json:"date_start"`
	DateEnd      time.Time     `json:"date_end"`
	PhotoURL     string        `json:"photo_url"`
	BgColor      []uint8       `json:"bg_color"`
	FgColor      []uint8       `json:"fg_color"`
}

type Events []Event
