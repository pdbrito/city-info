package cityinfoservice

import (
	"fmt"
	"github.com/pdbrito/city-info/internal/weatherservice"
	"github.com/pdbrito/city-info/internal/wikiservice"
	"github.com/sirupsen/logrus"
)

// CityInfoService defines our interface
type CityInfoService interface {
	CurrentStatus(cityName string) (CityInfo, error)
}

// CityInfo contains current information about a city.
type CityInfo struct {
	Description      string `json:"description"`
	WeatherSituation string `json:"weather_situation"`
	Temperature      string `json:"temperature"`
}

type cityInfoService struct {
	weatherService weatherservice.WeatherService
	wikiService    wikiservice.WikipediaService
	logger         logrus.Logger
}

// NewCityInfoService returns a new CityInfoService.
func NewCityInfoService(ws weatherservice.WeatherService, wis wikiservice.WikipediaService, logger logrus.Logger) CityInfoService {
	return cityInfoService{
		weatherService: ws,
		wikiService:    wis,
		logger:         logger,
	}
}

// CurrentStatus returns a CityInfo for the given city.
func (c cityInfoService) CurrentStatus(name string) (CityInfo, error) {
	weather, err := c.weatherService.CityWeather(name)
	if err != nil {
		err = fmt.Errorf("could not fetch weather for '%s': %v", name, err)
		c.logger.Error(err)
		return CityInfo{}, err
	}
	description, err := c.wikiService.CityDescription(name)
	if err != nil {
		err = fmt.Errorf("could not fetch description for '%s': %v", name, err)
		c.logger.Error(err)
		return CityInfo{}, err
	}

	return CityInfo{
		Description:      description,
		WeatherSituation: weather.Description,
		Temperature:      fmt.Sprintf("%.1f", weather.Temperature),
	}, nil
}
