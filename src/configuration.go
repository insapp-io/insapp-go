package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"log"
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

var mgoSession *mgo.Session

func Configuration() (Config, *mgo.DialInfo) {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)

	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatal("Error when parsing config file. Make sure the configuration file (config.json) is valid.")
	}

	info := &mgo.DialInfo{
		Addrs:    []string{"db"},
		Database: "insapp",
		Source:   "admin",
		Username: "insapp-admin",
		Password: configuration.DatabasePassword,
		Timeout:  time.Second * 10,
	}

	return configuration, info
}

//Creates a new session if mgoSession is nil i.e there is no active mongo session.
//If there is an active mongo session it will return a Clone
func GetMongoSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		_, info := Configuration()
		mgoSession, err = mgo.DialWithInfo(info)
		if err != nil {
			log.Fatal("Failed to start the Mongo session")
		}
		mgoSession.SetMode(mgo.Monotonic, true)
	}
	return mgoSession.Clone()
}

func (configuration Config) GetCDN() string {
	cdn := "https://insapp"
	if configuration.Environment == "dev" {
		cdn += ".insa-rennes"
	}
	cdn += ".fr/cdn/"

	return cdn
}
