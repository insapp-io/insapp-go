package main

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
)


func UploadNewImageController(w http.ResponseWriter, r *http.Request) {
	fileName := UploadImage(r)
	if fileName == "error" {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(bson.M{"error": "Failed to upload image"})
	} else {
    width, height := GetImageDimension(fileName)
    colors := GetImageColors(fileName)
    json.NewEncoder(w).Encode(bson.M{"file": fileName, "size": bson.M{"width": width, "height": height}, "colors": colors})
	}
}


func UploadImageController(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  fileName := UploadImageWithName(r, vars["name"])
	if fileName == "error" {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(bson.M{"error": "Failed to upload image"})
	} else {
		width, height := GetImageDimension(fileName)
    colors := GetImageColors(fileName)
		json.NewEncoder(w).Encode(bson.M{"file": fileName, "size": bson.M{"width": width, "height": height}, "colors": colors})
	}
}
