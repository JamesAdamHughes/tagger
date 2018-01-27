package handlers

import (
	"net/http"
	"fmt"
	"html/template"
	"golang.org/x/oauth2"
	"tagger/server/spotify"
	"github.com/zmb3/spotify"
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

func GetClientFromCookies(r *http.Request) *spotify.Client {
	c, _ := r.Cookie("authTokenTMP")
	if c != nil {
		fmt.Printf("%+v", c)

		// get user if possible
		client := spotify_manager.GetClient(&oauth2.Token{
			AccessToken: c.Value,
		})
		return client
	}
	return nil
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