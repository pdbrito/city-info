package http_test

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/pdbrito/city-info/internal/cityinfoservice"
	internalHttp "github.com/pdbrito/city-info/internal/http"
	"github.com/pdbrito/city-info/internal/weatherservice"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type weatherService struct {
	shouldError bool
	willReturn  weatherservice.Weather
}

func (w weatherService) CityWeather(name string) (weatherservice.Weather, error) {
	return w.willReturn, nil
}

func newWeatherService(shouldError bool, weather weatherservice.Weather) weatherService {
	return weatherService{
		shouldError: shouldError,
		willReturn:  weather,
	}
}

type wikiService struct {
	shouldError bool
	willReturn  string
}

func newWikiService(shouldError bool, description string) wikiService {
	return wikiService{
		shouldError: shouldError,
		willReturn:  description,
	}
}

func (w wikiService) CityDescription(name string) (string, error) {
	return w.willReturn, nil
}

func TestHandler_GetCityInfo(t *testing.T) {
	expectedTemp := 44.4
	cinfo := cityinfoservice.CityInfo{
		Description:      "Mos Eisley, you will never find a more wretched hive of scum and villainy.",
		WeatherSituation: "hot and sandy",
		Temperature:      fmt.Sprintf("%.1f", expectedTemp),
	}

	ws := newWeatherService(
		false,
		weatherservice.Weather{
			Temperature: expectedTemp,
			Description: cinfo.WeatherSituation,
		},
	)

	wis := newWikiService(
		false,
		cinfo.Description,
	)

	cis := cityinfoservice.NewCityInfoService(ws, wis, *logrus.New())

	handler := internalHttp.NewHandler(cis)
	server := httptest.NewServer(handler)

	queryURL := fmt.Sprintf(
		"%s/city-info?name=%s",
		server.URL,
		url.QueryEscape("Mos Eisley"),
	)

	res, err := http.Get(queryURL)

	if err != nil {
		t.Fatalf("err getting city info from handler; got %s, want nil", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf(
			"err getting city info from handler; got http status code %d, want %d",
			res.StatusCode,
			http.StatusOK,
		)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("err reading response body; got %s, want nil", err)
	}

	expectedJson, err := json.Marshal(cinfo)
	if err != nil {
		t.Fatalf("err marchalling city info into json; got %s, want nil", err)
	}

	if !cmp.Equal(body, expectedJson) {
		t.Fatalf(
			"err getting city info from handler; json did not match; diff:\n%s",
			cmp.Diff(
				string(body),
				string(expectedJson),
			),
		)
	}

}
