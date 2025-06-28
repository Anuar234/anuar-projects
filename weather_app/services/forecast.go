package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"weather_app/templates/cache"
)

type ForecastService struct {
	cache *cache.MemoryCache
}

type ForecastItem struct {
	Date        string
	Temp        float64
	Description string
	IconURL     string
}

type ForecastData struct {
	City      string
	Forecasts []ForecastItem
}

// initiliazes the forecast service with a cache instance
func NewForestSerivce(c *cache.MemoryCache) *ForecastService {
	return &ForecastService{cache: c}
}

func (s *ForecastService) GetForecast(city string) (ForecastData, error) {

	// normalize city name for cache
	key := fmt.Sprintf("%s_forecast", city)

	//chechk cache
	if data, ok := s.cache.Get(key); ok {
		return data.(ForecastData), nil
	}

	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		return ForecastData{}, fmt.Errorf("key not set")
	}
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/forecast?q=%s&appid=%s&units=metric&lang=en",
		city, apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return ForecastData{}, err
	}
	defer resp.Body.Close()

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return ForecastData{}, err
	}

	list, ok := raw["list"].([]interface{})
	if !ok {
		return ForecastData{}, fmt.Errorf("no forecast data available")
	}

	var forecasts []ForecastItem
	seenDays := map[string]bool{}

	for _, entry := range list {
		item := entry.(map[string]interface{})
		dtTxt := item["dt_txt"].(string)

		parsedTime, _ := time.Parse("2006-01-02 15:04:05", dtTxt)
		day := parsedTime.Weekday().String()

		if seenDays[day] {
			continue
		}
		seenDays[day] = true

		main := item["main"].(map[string]interface{})
		weather := item["weather"].([]interface{})[0].(map[string]interface{})

		forecasts = append(forecasts, ForecastItem{
			Date:        day,
			Temp:        main["temp"].(float64),
			Description: weather["description"].(string),
			IconURL:     fmt.Sprintf("https://openweathermap.org/img/wn/%s@2x.png", weather["icon"].(string)),
		})

		if len(forecasts) == 5 {
			break
		}
	}

	forecast := ForecastData{
		City:      city,
		Forecasts: forecasts,
	}

	s.cache.Set(key, forecast)

	return forecast, nil
}
