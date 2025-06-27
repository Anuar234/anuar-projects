package handlers

import (
	"html/template"
	"net/http"
	"weather_app/services"
)

func ForecastHandler(w http.ResponseWriter, r *http.Request) {
    city := r.URL.Query().Get("city")
    if city == "" {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    data, err := services.GetForecast(city)
    if err != nil {
        http.Error(w, "Error fetching forecast", http.StatusInternalServerError)
        return
    }

    addToRecent(city)

    tmpl := template.Must(template.ParseFiles("templates/forecast.html"))
    tmpl.Execute(w, data)
}
