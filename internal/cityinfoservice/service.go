package cityinfoservice

import (
	"fmt"
	"github.com/pdbrito/city-info/internal/weatherservice"
	"github.com/pdbrito/city-info/internal/wikiservice"
	"github.com/sirupsen/logrus"
	"sync"
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
	var wg sync.WaitGroup
	wg.Add(2)

	wChan, wsErr := make(chan weatherservice.Weather, 1), make(chan error, 1)
	go func() {
		defer wg.Done()
		weather, err := c.weatherService.CityWeather(name)
		if err != nil {
			err = fmt.Errorf("could not fetch weather for '%s': %v", name, err)
			c.logger.Error(err)
			wsErr <- err
			return
		}
		wChan <- weather
	}()

	dChan, wisErr := make(chan string, 1), make(chan error, 1)
	go func() {
		defer wg.Done()
		description, err := c.wikiService.CityDescription(name)
		if err != nil {
			err = fmt.Errorf("could not fetch descrpition for '%s': %v", name, err)
			c.logger.Error(err)
			wisErr <- err
			return
		}
		dChan <- description
	}()

	wg.Wait()

	select {
	case err := <-wsErr:
		return CityInfo{}, err
	case err := <-wisErr:
		return CityInfo{}, err
	default:
		w := <-wChan
		return CityInfo{
			Description:      <-dChan,
			WeatherSituation: w.Description,
			Temperature:      fmt.Sprintf("%.1f", w.Temperature),
		}, nil
	}
}
