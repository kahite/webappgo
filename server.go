package main

import (
	"fmt"
	"net/http"
	"log"
    "io/ioutil"
    "html/template"
)

type Page struct {
    Title string
    Body []byte
}

func (p *Page) save() error {
    filename := p.Title + ".txt"
    return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    return &Page{Title: title, Body: body}, nil
}

func renderTemplate(filename string, w http.ResponseWriter, p *Page) {
    t, _ := template.ParseFiles(filename + ".html")
    t.Execute(w, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
    pageTitle := r.URL.Path[len("/view/"):]
    p, err := loadPage(pageTitle)
    if err != nil {
        http.Redirect(w, r, "/edit/" + pageTitle, http.StatusNotFound)
        return 
    }
    renderTemplate("view", w, p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	pageTitle := r.URL.Path[len("/edit/"):]
    p, err := loadPage(pageTitle)
    if err != nil {
        p = &Page{Title: pageTitle}
    }
    renderTemplate("edit", w, p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	pageTitle := r.URL.Path[len("/save/"):]
	p, _ := loadPage(pageTitle)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", editHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
