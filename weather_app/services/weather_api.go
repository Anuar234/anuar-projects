package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"weather_app/templates/cache"
)

type WeatherAPIService struct {
	cache *cache.MemoryCache
}

type WeatherData struct {
	City        string
	Temp        float64
	Description string
	IconURL     string
}

// Initializes the weather service with a cahce in
func NewWeatherAPIService(c *cache.MemoryCache) *WeatherAPIService {
	return &WeatherAPIService{cache: c}
}

// Пример функции
func (s *WeatherAPIService) GetCurrentWeather(city string) (WeatherData, error) {

	//normalize city name for cache key
	key := fmt.Sprintf("%s_current", city)

	//chechking the cache
	if data, ok := s.cache.Get(key); ok {
		return data.(WeatherData), nil
	}

	//cache miss, fetch from API
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		return WeatherData{}, fmt.Errorf("Key not set in env")
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return WeatherData{}, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return WeatherData{}, err
	}

	weather := WeatherData{
		City:        city,
		Temp:        data["main"].(map[string]interface{})["temp"].(float64),
		Description: data["weather"].([]interface{})[0].(map[string]interface{})["description"].(string),
		IconURL:     fmt.Sprintf("https://openweathermap.org/img/wn/%s@2x.png", data["weather"].([]interface{})[0].(map[string]interface{})["icon"].(string)),
	}

	// store in cache
	s.cache.Set(key, weather)

	return weather, nil
}
