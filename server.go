package main

import (
	"net/http"
	"log"
    "io/ioutil"
    "html/template"
    "regexp"
)

/*** Globals ***/
var templates = template.Must(template.ParseGlob("templ/*.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

/*** Structs ***/
type Page struct {
    Title string
    Body []byte
}

/*** Utils ***/
func loadPage(title string) (*Page, error) {
    filename := "data/" + title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    return &Page{Title: title, Body: body}, nil
}

func (p *Page) save() error {
    filename := "data/" + p.Title + ".txt"
    return ioutil.WriteFile(filename, p.Body, 0600)
}

func renderTemplate(filename string, w http.ResponseWriter, p *Page) {
    err := templates.ExecuteTemplate(w, filename + ".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

/*** Route handlers ***/
func viewHandler(w http.ResponseWriter, r *http.Request, pageTitle string) {
    p, err := loadPage(pageTitle)
    if err != nil {
        http.Redirect(w, r, "/edit/" + pageTitle, http.StatusNotFound)
        return 
    }
    renderTemplate("view", w, p)
}

func editHandler(w http.ResponseWriter, r *http.Request, pageTitle string) {
    p, err := loadPage(pageTitle)
    if err != nil {
        p = &Page{Title: pageTitle}
    }
    renderTemplate("edit", w, p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, pageTitle string) {
    body := r.FormValue("body")
    p := &Page{Title: pageTitle, Body: []byte(body)}
    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/" + pageTitle, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
        }
		fn(w, r, m[2])
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
}

/*** Main function ***/
func main() {
    http.HandleFunc("/", indexHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
