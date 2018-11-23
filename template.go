package main

import "html/template"

var indexTmpl = template.Must(template.ParseFiles("layout.html"))
