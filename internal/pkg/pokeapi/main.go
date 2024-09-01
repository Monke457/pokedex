package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	Next *string 
	Previous *string 
}

type JsonLocations struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocations(url *string) (*JsonLocations, error) {
	var endpoint string
	if url == nil {
		endpoint = "https://pokeapi.co/api/v2/location-area/"
	} else {
		endpoint = *url
	}

	res, err := get(endpoint)
	if err != nil {
		return nil, err
	}

	jsonLocs := JsonLocations{}
	err = json.Unmarshal(res, &jsonLocs)
	if err != nil {
		return nil, err
	}

	return &jsonLocs, nil
}

func get(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return []byte{}, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}
