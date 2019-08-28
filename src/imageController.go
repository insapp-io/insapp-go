package main

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
)

// UploadNewImageController will upload a new image in the cdn
func UploadNewImageController(w http.ResponseWriter, r *http.Request) {
	fileName, err := UploadImage(r)
	ResponseHandler(&w, fileName, err)
}

// UploadImageController will upload a new image in the cdn
func UploadImageController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName, err := UploadImageWithName(r, vars["name"])
	ResponseHandler(&w, fileName, err)
}

// ResponseHandler will response to the client
func ResponseHandler(w *http.ResponseWriter, fileName string, err error) {
	if err != nil || fileName == "" {
		(*w).WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(*w).Encode(bson.M{"error": "Failed to upload image"})
	} else {
		width, height := GetImageDimension(fileName)
		if width == 0 || height == 0 {
			_ = json.NewEncoder(*w).Encode(bson.M{"error": "Bad image format"})
			return
		}
		colors := GetImageColors(fileName)
		_ = json.NewEncoder(*w).Encode(bson.M{"file": fileName, "size": bson.M{"width": width, "height": height}, "colors": colors})
	}
}
