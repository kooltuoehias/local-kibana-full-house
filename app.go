package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Page struct {
	Body []byte
}

var templates = template.Must(template.ParseFiles("app.html"))
var buffer bytes.Buffer

func handler(w http.ResponseWriter, r *http.Request) {
	if buffer.Len() == 0 {
		p := &Page{}
		renderTemplate(w, "app", p)
	} else {
		p := &Page{Body: buffer.Bytes()}
		renderTemplate(w, "app", p)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	services := r.Form["services"]
	env := r.FormValue("envs")
	fromDate := r.FormValue("fromDate")
	toDate := r.FormValue("toDate")
	buffer.Reset()
	buffer.WriteString(strings.Join(services, ", ") + "\n" + env + "\n" + fromDate + " <=====> " + toDate + "\n")
	for _, service := range services {
		Do(env, service, &buffer)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

func main() {
	http.HandleFunc("/", makeHandler(handler))
	http.HandleFunc("/save", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":9040", nil))
}
