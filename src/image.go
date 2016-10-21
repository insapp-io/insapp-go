package main

import (
    "io"
    "math/rand"
    "net/http"
    "os"
  	"time"
    "os/exec"

  	"image"
  	_ "image/jpeg"
  	_ "image/png"

    "strings"
    "strconv"

)

func UploadImage(r *http.Request) string{
  return UploadImageWithName(r, RandomString(40))
}

// UploadImage will manage the upload image from a POST request
func UploadImageWithName(r *http.Request, name string) string {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		return "error"
	}
	defer file.Close()

	fileName := name
	f, err := os.OpenFile("./img/"+fileName+".png", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "error"
	}

	defer f.Close()
	io.Copy(f, file)

	return fileName + ".png"
}

func GetImageDimension(fileName string) (int, int) {
    file, _ := os.Open("./img/"+fileName)
    image, _, _ := image.DecodeConfig(file)
    return image.Width, image.Height
}

func GetImageColors(fileName string) [][]int {
  var result [][]int

  bytes, err := exec.Command("python", "color-thief.py", "./img/" + fileName).Output()

  if err != nil {
    return result
  }

  out := string(bytes)
  out = strings.Replace(out, "[", "", -1)
  out = strings.Replace(out, "]", "", -1)
  out = strings.Replace(out, " ", "", -1)
  out = strings.Replace(out, ",", " ", -1)
  out = strings.Replace(out, ")", "", -1)
  out = strings.Replace(out, "(", "", 1)
  split := strings.Split(out, "(")

  for _, colorData := range split {

    var colors []int

    colorData = strings.Replace(colorData, "(", "", -1)
    colorData = strings.Replace(colorData, ")", "", -1)
    stringColors := strings.Split(colorData, " ")

    for _, col := range stringColors {
      i, err := strconv.Atoi(strings.TrimSpace(col))
      if err == nil {
        colors = append(colors, i)
      }
    }
    result = append(result, colors)
  }
  return result
}

// RandomString generates a randomString (y)
func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
