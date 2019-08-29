package main

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	resize "github.com/nfnt/resize"
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

// UploadImage will give a random name and upload the image
func UploadImage(r *http.Request) (string, error) {
	return UploadImageWithName(r, RandomString(40))
}

// UploadImageWithName will manage the upload image from a POST request
func UploadImageWithName(r *http.Request, name string) (string, error) {
	_ = r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	imgType := mimeFromIncipit(data)

	if imgType == "" {
		return "", errors.New("Can't get image format")
	}
	defer file.Close()

	fileName := name
	err = ioutil.WriteFile("./img/"+fileName+"."+imgType, data, 0666)
	if err != nil {
		return "", err
	}
	return fileName + "." + imgType, nil
}

// GetImageDimension will return image dimention in pixels
func GetImageDimension(fileName string) (int, int) {
	file, err := os.Open("./img/" + fileName)
	if err != nil {
		return 0, 0
	}
	decodedImage, _, _ := image.DecodeConfig(file)
	return decodedImage.Width, decodedImage.Height
}

// GetImageColors will return a palette of colors found in the image
func GetImageColors(fileName string) [][]int {
	return extractor.NewExtractor("./img/"+fileName, 10).GetPalette(6)
}

// GetImagesNames will return a string of all images in cdn
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

// ResizeImage will resize an image as jpeg and response its new name
// Put 0 to newWidth or newHeight to make it automatically calculate to keep aspect ratio.
func ResizeImage(imageName string, newWidth uint, newHeight uint) (string, error) {
	file, err := os.Open("./img/" + imageName)
	defer file.Close()
	if err != nil {
		return "", err
	}

	origanialImage, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	newImg := image.NewRGBA(origanialImage.Bounds())
	// paste a white background
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	// paste original image
	draw.Draw(newImg, newImg.Bounds(), origanialImage, origanialImage.Bounds().Min, draw.Over)

	finalImage := resize.Resize(newWidth, newHeight, newImg, resize.Lanczos3)
	name := RandomString(40)
	out, err := os.Create("./img/" + name + ".jpeg")
	if err != nil {
		return "", err
	}
	defer out.Close()

	err = jpeg.Encode(out, finalImage, nil)
	return name + ".jpeg", err
}

// ArchiveImage will move images in a subdirectory "archive"
func ArchiveImage(fileName string) error {
	if _, err := os.Stat("./img/archive/"); os.IsNotExist(err) {
		os.Mkdir("./img/archive/", 0755) // rwxr-wr-x
	}
	oldLocation := "./img/" + fileName
	newLocation := "./img/archive/" + fileName
	err := os.Rename(oldLocation, newLocation)
	return err
}

// DeleteImage will permanently delete an image
func DeleteImage(fileName string) error {
	err := os.Remove("./img/" + fileName)
	return err
}
