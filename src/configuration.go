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
	DatabaseName     string `json:"mongo_database_name"`
	DatabaseSource   string `json:"mongo_database_source"`
	DatabaseUsername string `json:"mongo_database_username"`
	DatabasePassword string `json:"mongo_database_password"`
	Environment      string `json:"env"`
	Port             string `json:"port"`
}

var mgoSession *mgo.Session

func Configuration() (Config, *mgo.DialInfo) {
	file, err1 := os.Open("config.json")
	decoder := json.NewDecoder(file)

	if err1 != nil {
		log.Fatal(err1)
	}

	configuration := Config{}
	err2 := decoder.Decode(&configuration)
	if err2 != nil {
		log.Fatal("Error when parsing config file. Make sure the configuration file (config.json) is valid.")
	}

	var address []string
	if configuration.Environment == "local" {
		address = append(address, "localhost:27017")
	} else {
		address = append(address, "db")
	}

	info := &mgo.DialInfo{
		Addrs:    address,
		Database: configuration.DatabaseName,
		Source:   configuration.DatabaseSource,
		Username: configuration.DatabaseUsername,
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
			log.Fatal(err)
			log.Fatal("Failed to start the Mongo session")
		}
		mgoSession.SetMode(mgo.Monotonic, true)
	}
	return mgoSession.Clone()
}

func (configuration Config) GetCDN() string {
	var cdn string

	switch configuration.Environment {
	case "prod":
		cdn = "https://insapp.fr/cdn/"
	case "dev":
		cdn = "https://insapp.insa-rennes.fr/cdn/"
	case "local":
		cdn = "test"
	}

	return cdn
}
