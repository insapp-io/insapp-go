package insapp

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"gopkg.in/mgo.v2"
)

// Config defines how to model a Config
type Config struct {
	Domain           string `json:"domain"`
	Environment      string `json:"env"`
	GoogleEmail      string `json:"google_email"`
	GooglePassword   string `json:"google_password"`
	FirebaseKey      string `json:"firebase_key"`
	DatabaseName     string `json:"mongo_database_name"`
	DatabaseSource   string `json:"mongo_database_source"`
	DatabaseUsername string `json:"mongo_database_username"`
	DatabasePassword string `json:"mongo_database_password"`
	PrivateKeyPath   string `json:"private_key_path"`
	PublicKeyPath    string `json:"public_key_path"`
	Port             string `json:"port"`
}

var mgoSession *mgo.Session
var config *Config

// InitConfig loads the configuration from the filesystem.
func InitConfig() *Config {
	file, err1 := os.Open("config.json")
	decoder := json.NewDecoder(file)

	if err1 != nil {
		log.Fatal(err1)
	}

	err2 := decoder.Decode(&config)
	if err2 != nil {
		log.Fatal("Error when parsing config file. Make sure the configuration file (config.json) is valid.")
	}

	return config
}

// GetMongoSession creates a new session.
// If there is an active mongo session it will return a Clone.
func GetMongoSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.DialWithInfo(initMongoConfig())

		if err != nil {
			log.Fatal(err)
			log.Fatal("Failed to start the Mongo session")
		}

		mgoSession.SetMode(mgo.Monotonic, true)
	}

	return mgoSession.Clone()
}

// GetCDN returns the CDN address depending the configuration.
func (config Config) GetCDN() string {
	var cdn string

	switch config.Environment {
	case "prod":
		cdn = "https://" + config.Domain + "/cdn/"
	case "dev":
		cdn = "https://" + config.Domain + "/cdn/"
	case "local":
		cdn = "test"
	}

	return cdn
}

func initMongoConfig() *mgo.DialInfo {
	var address []string
	if config.Environment == "local" {
		address = append(address, "localhost:27017")
	} else {
		address = append(address, "db")
	}

	return &mgo.DialInfo{
		Addrs:    address,
		Database: config.DatabaseName,
		Source:   config.DatabaseSource,
		Username: config.DatabaseUsername,
		Password: config.DatabasePassword,
		Timeout:  time.Second * 10,
	}
}
