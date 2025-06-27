package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type WeatherData struct {
    City        string
    Temp        float64
    Description string
    IconURL     string
}

// Пример функции
func GetCurrentWeather(city string) (WeatherData, error) {
    apiKey := os.Getenv("OPENWEATHER_API_KEY")
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

    return WeatherData{
        City:        city,
        Temp:        data["main"].(map[string]interface{})["temp"].(float64),
        Description: data["weather"].([]interface{})[0].(map[string]interface{})["description"].(string),
        IconURL:     fmt.Sprintf("https://openweathermap.org/img/wn/%s@2x.png", data["weather"].([]interface{})[0].(map[string]interface{})["icon"].(string)),
    }, nil
}
