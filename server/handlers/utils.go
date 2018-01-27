package handlers

import (
	"net/http"
	"fmt"
	"html/template"
	"io/ioutil"
)

const STATIC_PATH = "server/static"

type Page struct {
	Title string
	Body  string
}

type AuthPage struct {
	Title string
	Body string
	AuthUrl string
}

func RenderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	filename := STATIC_PATH + "/html/" + tmpl + ".html"

	t, err := template.ParseFiles(filename)
	if err != nil {
		fmt.Print(err)
		MissingPageHandler(w, filename)
		return
	}

	t.Execute(w, p)
}

func MissingPageHandler(w http.ResponseWriter, filename string) {
	t, _ := template.ParseFiles(STATIC_PATH + "/html/error_404.html")
	t.Execute(w, Page{
		Title: "File Not Found",
		Body: filename,
	})
}

func LoadPage(title string) (*Page, error) {
	filename := STATIC_PATH + "/html/" + title + ".html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	return &Page{Title: title, Body: string(body)}, nil
}