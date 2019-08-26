package main

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	extractor "github.com/thomas-bouvier/palette-extractor"
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
	_ = r.ParseMultipartForm(32 << 20)
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
	decodedImage, _, _ := image.DecodeConfig(file)
	return decodedImage.Width, decodedImage.Height
}

func GetImageColors(fileName string) [][]int {
	return extractor.NewExtractor("./img/"+fileName, 10).GetPalette(6)
}

func GetImagesNames() ([]string, error) {
	files, err := ioutil.ReadDir("./img/")
	if err != nil {
		return nil, err
	}

	var result []string
	for _, file := range files {
		if file.Name() != ".gitignore" && file.Name() != "archive" && file.Name() != "index.html" {
			result = append(result, file.Name())
		}
	}
	return result, nil
}

func ArchiveImage(fileName string) error {
	if _, err := os.Stat("./img/archive/"); os.IsNotExist(err) {
		os.Mkdir("./img/archive/", 0755) // rwxr-wr-x
	}
	oldLocation := "./img/" + fileName
	newLocation := "./img/archive/" + fileName
	err := os.Rename(oldLocation, newLocation)
	return err
}

func DeleteImage(fileName string) error {
	err := os.Remove("./img/" + fileName)
	return err
}
