package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const STATIC_PATH = "/src/tagger/static"

type Page struct {
	Title string
	Body  string
}

type AuthPage struct {
	Title   string
	Body    string
	AuthUrl string
}

func RenderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	filename := getTemplateFileName(tmpl)
	t, err := template.ParseFiles(filename)
	if err != nil {
		fmt.Println(err.Error())
		MissingPageHandler(w, tmpl)
		return
	}

	_ = t.Execute(w, p)
}

func MissingPageHandler(w http.ResponseWriter, missingPage string) {
	filename := getTemplateFileName("error_404")
	t, _ := template.ParseFiles(filename)
	_ = t.Execute(w, Page{
		Title: "File Not Found",
		Body:  missingPage,
	})
}

func LoadTemplate(title string) (*Page, error) {
	filename := getTemplateFileName(title)
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	return &Page{Title: title, Body: string(body)}, nil
}

func getTemplateFileName(tmpl string) string {
	d, _ := os.Getwd()
	filename := filepath.Join(d, STATIC_PATH, "/html/"+tmpl+".html")
	fmt.Println(filename)
	return filename
}

func returnJson(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(i)
	return
}

func returnError(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, fmt.Sprintf("%+v", i), http.StatusInternalServerError)
	return
}
