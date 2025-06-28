package handlers

import (
	"html/template"
	"net/http"
	"weather_app/services"
)

type WeatherHadnler struct {
	service *services.WeatherAPIService
}

func NewWeatherHandler(s *services.WeatherAPIService) *WeatherHadnler {
	return &WeatherHadnler{service: s}
}

func (h *WeatherHadnler) WeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data, err := h.service.GetCurrentWeather(city)
	if err != nil {
		http.Error(w, "Error fetching weather", http.StatusInternalServerError)
		return
	}

	addToRecent(city)

	tmpl := template.Must(template.ParseFiles("templates/weather.html"))
	tmpl.Execute(w, data)
}
