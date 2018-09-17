package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"os"
	"time"
)

// Post defines how to model a Post
type Config struct {
	GoogleEmail      string `json:"google_email"`
	GooglePassword   string `json:"google_password"`
	FirebaseKey      string `json:"firebase_key"`
	DatabasePassword string `json:"mongo_database_password"`
	Environment      string `json:"env"`
	Port             string `json:"port"`
}

func Configuration() (Config, *mgo.DialInfo, error) {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)

	configuration := Config{}
	err := decoder.Decode(&configuration)

	info := &mgo.DialInfo{
		Addrs:    []string{"db"},
		Database: "insapp",
		Source:   "admin",
		Username: "insapp-admin",
		Password: configuration.DatabasePassword,
		Timeout:  time.Second * 10,
	}

	return configuration, info, err
}

func (configuration Config) GetCDN() string {
	cdn := "https://insapp"
	if configuration.Environment == "dev" {
		cdn += ".insa-rennes"
	}
	cdn += ".fr/cdn/"

	return cdn
}
