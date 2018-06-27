package main

import (
	"github.com/thomas-bouvier/palette-extractor"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var magicTable = map[string]string{
	"\xff\xd8\xff":      "jpeg",
	"\x89PNG\r\n\x1a\n": "png",
	"GIF87a":            "gif",
	"GIF89a":            "gif",
}

func mimeFromIncipit(incipit []byte) string {
	incipitStr := string(incipit)
	for magic, mime := range magicTable {
		if strings.HasPrefix(incipitStr, magic) {
			return mime
		}
	}
	return ""
}

func UploadImage(r *http.Request) string {
	return UploadImageWithName(r, RandomString(40))
}

// UploadImage will manage the upload image from a POST request
func UploadImageWithName(r *http.Request, name string) string {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		return "error"
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "error"
	}
	imgType := mimeFromIncipit(data)

	if imgType == "" {
		return "error"
	}
	defer file.Close()

	fileName := name
	err = ioutil.WriteFile("./img/"+fileName+"."+imgType, data, 0666)
	if err != nil {
		return "error"
	}
	return fileName + "." + imgType
}

func GetImageDimension(fileName string) (int, int) {
	file, err := os.Open("./img/" + fileName)
	if err != nil {
		return 0, 0
	}
	image, _, _ := image.DecodeConfig(file)
	return image.Width, image.Height
}

func GetImageColors(fileName string) [][]int {
	return extractor.NewExtractor("./img/"+fileName, 10).GetPalette(6)
}
