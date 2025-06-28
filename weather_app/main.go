package main

import (
	"fmt"
	"net/http"
	"time"

	"weather_app/handlers"
	"weather_app/services"
	"weather_app/templates/cache"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Could not load .env")
	}

	// initialize cache with 10-minute ttl
	cacheInstance := cache.NewMemoryCache(10 * time.Minute)

	// start cleanup goroutine
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			cacheInstance.Cleanup()
		}
	}()

	weatherService := services.NewWeatherAPIService(cacheInstance)
	forecastService := services.NewForestSerivce(cacheInstance)

	// Initialize handlers with services
	weatherHandler := handlers.NewWeatherHandler(weatherService)
	forecastHandler := handlers.NewForecastHandler(forecastService)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/search", handlers.SearchHandler)
	http.HandleFunc("/weather", weatherHandler.WeatherHandler)
	http.HandleFunc("/forecast", forecastHandler.ForecastHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
