package main

import (
	"fmt"
	"net/http"

	"weather_app/handlers"

	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        fmt.Println("Could not load .env")
    }

    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/search", handlers.SearchHandler)
    http.HandleFunc("/weather", handlers.WeatherHandler)
    http.HandleFunc("/forecast", handlers.ForecastHandler)

    fmt.Println("Server running on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
