package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func loadPage(title string) ([]byte, error) {
	filename := "pages/" + title + ".html"
	page, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return page, nil
}

// Index is just a test actually
func Index(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "Insapp REST API - v1.0")
}

func HowToPost(w http.ResponseWriter, r *http.Request) {
	p, _ := loadPage("how-to-post")
	fmt.Fprintf(w, "%s", p)
}

func Credit(w http.ResponseWriter, r *http.Request) {
	p, _ := loadPage("credit")
	fmt.Fprintf(w, "%s", p)
}

func Legal(w http.ResponseWriter, r *http.Request) {
	p, _ := loadPage("legal")
	fmt.Fprintf(w, "%s", p)
}
