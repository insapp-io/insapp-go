package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"os"
)

// Post defines how to model a Post
type Config struct {
	GoogleEmail      string `json:"google_email"`
	GooglePassword   string `json:"google_password"`
	GoogleKey        string `json:"google_key"`
	DatabaseName     string `json:"mongo_database_name"`
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
		Database: configuration.DatabaseName,
		Username: "Insapp",
		Password: configuration.DatabasePassword,
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
