package main

import (
    "encoding/json"
    "os"
)

// Post defines how to model a Post
type Config struct {
	Email       string      `json:"email"`
	Password    string      `json:"password"`
	GoogleKey   string      `json:"googlekey"`
  Environment string      `json:"env"`
}


func Configuration() (Config, error) {
  file, _ := os.Open("config.json")
  decoder := json.NewDecoder(file)
  configuration := Config{}
  err := decoder.Decode(&configuration)
  return configuration, err
}
