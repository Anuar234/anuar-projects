package handlers

import (
	"html/template"
	"net/http"
)

var RecentCities []string

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("templates/index.html"))
    tmpl.Execute(w, RecentCities)
}
