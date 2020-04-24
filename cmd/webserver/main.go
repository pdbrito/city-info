package main

import (
	"github.com/pdbrito/city-info/internal/cityinfoservice"
	httpTransport "github.com/pdbrito/city-info/internal/http"
	"github.com/pdbrito/city-info/internal/weatherservice"
	"github.com/pdbrito/city-info/internal/wikiservice"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
)

func main() {
	// Used to authenticate requests to the OWS API
	apiKey := os.Getenv("OWS_API_KEY")

	if apiKey == "" {
		log.Fatal("OWS_API_KEY environment variable not set, exiting")
	}

	weatherService := weatherservice.NewWeatherService(apiKey)
	wikipediaService := wikiservice.NewWikipediaService()

	cityInfoService := cityinfoservice.NewCityInfoService(weatherService, wikipediaService, *logrus.New())
	cityInfoHandler := httpTransport.NewHandler(cityInfoService)

	if err := http.ListenAndServe("localhost:8181", cityInfoHandler); err != nil {
		log.Fatalf("server listen error: %+v", err)
	}
}
