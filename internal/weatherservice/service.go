package weatherservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// Weather contains information about the current weather.
type Weather struct {
	Temperature float64
	Description string
}

// WeatherService defines an interface for a wikipedia service
type WeatherService interface {
	CityWeather(name string) (Weather, error)
}

const endpoint string = "https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s"

type weatherService struct {
	apiKey string
}

// NewWeatherService returns a new WeatherService.
func NewWeatherService(apiKey string) WeatherService {
	return weatherService{
		apiKey: apiKey,
	}
}

// CityWeather returns the current weather for a given city from openWeatherMap.
func (w weatherService) CityWeather(name string) (Weather, error) {
	if name == "" {
		return Weather{}, errors.New("invalid name; should be non empty string")
	}

	resp, err := http.Get(fmt.Sprintf(endpoint, url.QueryEscape(name), w.apiKey))
	if err != nil {
		return Weather{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Weather{}, errors.New("city not found")
	}

	var apiResp apiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return Weather{}, err
	}

	var forecast Weather
	forecast.Temperature = apiResp.Main.Temp

	for _, w := range apiResp.Weather {
		forecast.Description = w.Description
		break // only interested in first result
	}

	return forecast, nil
}

type apiResponse struct {
	Base   string `json:"base"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Cod   int `json:"cod"`
	Coord struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"coord"`
	Dt   int `json:"dt"`
	ID   int `json:"id"`
	Main struct {
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
		Pressure  int     `json:"pressure"`
		Temp      float64 `json:"temp"`
		TempMax   float64 `json:"temp_max"`
		TempMin   float64 `json:"temp_min"`
	} `json:"main"`
	Name string `json:"name"`
	Sys  struct {
		Country string `json:"country"`
		ID      int    `json:"id"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
		Type    int    `json:"type"`
	} `json:"sys"`
	Timezone   int `json:"timezone"`
	Visibility int `json:"visibility"`
	Weather    []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
		ID          int    `json:"id"`
		Main        string `json:"main"`
	} `json:"weather"`
	Wind struct {
		Deg   int     `json:"deg"`
		Speed float64 `json:"speed"`
	} `json:"wind"`
}
