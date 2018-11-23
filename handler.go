package main

import (
	"html/template"
	"log"
	"net/http"
)

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	indexTmpl, err := template.ParseFiles("html/layout.html", "html/index.html", "html/navbar.html")

	err = indexTmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		http.Error(w, "blog: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func registerPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	registerTmpl, err := template.ParseFiles("html/layout.html", "html/register.html", "html/navbar.html")

	err = registerTmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		http.Error(w, "blog: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func startServer() error {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/":
			indexPageHandler(w, r)
		case r.Method == http.MethodGet && r.URL.Path == "/register/":
			registerPageHandler(w, r)
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
