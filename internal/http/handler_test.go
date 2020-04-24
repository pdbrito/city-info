package http_test

import (
	"encoding/json"
	"errors"
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
	"strings"
	"testing"
)

type weatherService struct {
	err        error
	willReturn weatherservice.Weather
}

func (w weatherService) CityWeather(name string) (weatherservice.Weather, error) {
	return w.willReturn, w.err
}

func newWeatherService(err error, weather weatherservice.Weather) weatherService {
	return weatherService{
		err:        err,
		willReturn: weather,
	}
}

type wikiService struct {
	err        error
	willReturn string
}

func newWikiService(err error, description string) wikiService {
	return wikiService{
		err:        err,
		willReturn: description,
	}
}

func (w wikiService) CityDescription(name string) (string, error) {
	return w.willReturn, w.err
}

func TestHandler_GetCityInfo(t *testing.T) {
	expectedTemp := 44.4
	cinfo := cityinfoservice.CityInfo{
		Description:      "Mos Eisley, you will never find a more wretched hive of scum and villainy.",
		WeatherSituation: "hot and sandy",
		Temperature:      fmt.Sprintf("%.1f", expectedTemp),
	}

	ws := newWeatherService(
		nil,
		weatherservice.Weather{
			Temperature: expectedTemp,
			Description: cinfo.WeatherSituation,
		},
	)

	wis := newWikiService(
		nil,
		cinfo.Description,
	)

	logger := logrus.New()
	logger.Out = ioutil.Discard
	cis := cityinfoservice.NewCityInfoService(ws, wis, *logger)

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

func TestHandler_GetCityInfo_ReturnsErrorResponseWhenWeatherServiceErrors(t *testing.T) {
	wsErr := "something went wrong"
	ws := newWeatherService(
		errors.New(wsErr),
		weatherservice.Weather{},
	)

	wis := newWikiService(
		nil,
		"Mos Eisley, you will never find a more wretched hive of scum and villainy.",
	)

	logger := logrus.New()
	logger.Out = ioutil.Discard
	cis := cityinfoservice.NewCityInfoService(ws, wis, *logger)

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

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf(
			"err getting city info from handler; got http status code %d, want %d",
			res.StatusCode,
			http.StatusNotFound,
		)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("err reading response body; got %v, want nil", err)
	}

	var errResp internalHttp.ErrorResponse
	err = json.Unmarshal(body, &errResp)
	if err != nil {
		t.Fatalf("err unmarshalling json response; got %v want nil", err)
	}

	if !strings.Contains(errResp.Error, wsErr) {
		t.Fatalf(
			"unexpected error message in error response; got %v, want %v",
			errResp.Error,
			wsErr,
		)
	}
}

func TestHandler_GetCityInfo_ReturnsErrorResponseWhenWikiServiceErrors(t *testing.T) {
	ws := newWeatherService(
		nil,
		weatherservice.Weather{
			Temperature: 44.4,
			Description: "hot and sandy",
		},
	)

	wisErr := "something went wrong"
	wis := newWikiService(
		errors.New(wisErr),
		"",
	)

	logger := logrus.New()
	logger.Out = ioutil.Discard
	cis := cityinfoservice.NewCityInfoService(ws, wis, *logger)

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

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf(
			"err getting city info from handler; got http status code %d, want %d",
			res.StatusCode,
			http.StatusNotFound,
		)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("err reading response body; got %v, want nil", err)
	}

	var errResp internalHttp.ErrorResponse
	err = json.Unmarshal(body, &errResp)
	if err != nil {
		t.Fatalf("err unmarshalling json response; got %v want nil", err)
	}

	if !strings.Contains(errResp.Error, wisErr) {
		t.Fatalf(
			"unexpected error message in error response; got %v, want %v",
			errResp.Error,
			wisErr,
		)
	}
}
