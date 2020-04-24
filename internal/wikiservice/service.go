package wikiservice

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// WikipediaService defines an interface for a wikipedia service
type WikipediaService interface {
	CityDescription(name string) (string, error)
}

const endpoint string = "https://en.wikipedia.org/w/api.php?action=query&prop=extracts&exsentences=1&exintro&explaintext&format=json&redirects&titles="

type wikipediaService struct {
}

// NewWikipediaService returns a new WikipediaService.
func NewWikipediaService() WikipediaService {
	return wikipediaService{}
}

type page struct {
	Pageid  int    `json:"pageid"`
	Ns      int    `json:"ns"`
	Title   string `json:"title"`
	Extract string `json:"extract"`
}

type apiResponse struct {
	Batchcomplete string `json:"batchcomplete"`
	Query         struct {
		Pages map[string]page `json:"pages"`
	} `json:"query"`
}

// CityDescription returns a description of the passed in city from wikipedia.
func (w wikipediaService) CityDescription(name string) (string, error) {
	if name == "" {
		return "", errors.New("invalid name; should be non empty string")
	}

	resp, err := http.Get(endpoint + url.QueryEscape(name))

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var apiResp apiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return "", err
	}

	if len(apiResp.Query.Pages) < 1 {
		return "", errors.New("empty response from wikipedia api")
	}

	cityDescriptions := make([]string, 0)
	for _, v := range apiResp.Query.Pages {
		cityDescriptions = append(cityDescriptions, v.Extract)
		break // only interested in the first result
	}

	if cityDescriptions[0] == "" {
		return "", errors.New("empty response from wikipedia api")
	}

	return cityDescriptions[0], nil
}
