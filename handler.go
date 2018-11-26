package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	indexTmpl, err := template.ParseFiles("html/layout.html", "html/index.html", "html/navbar.html", "html/footer.html")

	err = indexTmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		http.Error(w, "blog: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	registerTmpl, err := template.ParseFiles("html/layout.html", "html/eventManagement.html", "html/navbar.html", "html/footer.html")

	err = registerTmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		http.Error(w, "blog: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func eventAddingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	registerTmpl, err := template.ParseFiles("html/layout.html", "html/addEvent.html", "html/navbar.html", "html/footer.html")

	err = registerTmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		http.Error(w, "blog: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func addDataToDBHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "here")

}

func startServer() error {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/":
			indexPageHandler(w, r)
		case r.Method == http.MethodGet && r.URL.Path == "/eventmanagement/":
			eventHandler(w, r)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/eventmanagement/") && strings.HasSuffix(r.URL.Path, "/add/"):
			eventAddingHandler(w, r)
		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/eventmanagement/") && strings.HasSuffix(r.URL.Path, "/add/toDB/"):
			addDataToDBHandler(w, r)
		}
	})

	http.Handle("/css/", http.FileServer(http.Dir("./html/assets/")))
	http.Handle("/fonts/", http.FileServer(http.Dir("./html/assets/")))
	http.Handle("/img/", http.FileServer(http.Dir("./html/assets/")))
	http.Handle("/js/", http.FileServer(http.Dir("./html/assets/")))
	http.Handle("/scss/", http.FileServer(http.Dir("./html/assets/")))
	http.Handle("/vendor/", http.FileServer(http.Dir("./html/assets/")))

	return http.ListenAndServe(":8000", nil)
}
