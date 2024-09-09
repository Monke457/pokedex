package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	cache "pokedex/internal/pkg/pokecache"
)

var c = cache.NewCache(5 * time.Second)

type Config struct {
	Next *string 
	Previous *string 
	Args []string
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

type JsonPokes struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func GetLocations(url *string) (*JsonLocations, error) {
	var endpoint string
	if url == nil {
		endpoint = "https://pokeapi.co/api/v2/location-area/"
	} else {
		endpoint = *url
	}

	var err error
	res, ok := c.Get(endpoint)

	if !ok {
		res, err = get(endpoint)
		if err != nil {
			return nil, err
		}
		c.Add(endpoint, res)
	}

	jsonLocs := JsonLocations{}
	err = json.Unmarshal(res, &jsonLocs)
	if err != nil {
		return nil, err
	}
	return &jsonLocs, nil
}

func ExploreLocation(loc string) ([]string, error) {
	endpoint := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", loc)

	var err error
	res, ok := c.Get(endpoint)

	if !ok {
		res, err = get(endpoint)
		if err != nil {
			return nil, err
		}
		c.Add(endpoint, res)
	}

	jsonPokes := JsonPokes{}
	err = json.Unmarshal(res, &jsonPokes)
	if err != nil {
		return nil, err
	}

	pokemon := []string{}
	for _, p := range jsonPokes.PokemonEncounters {
		pokemon = append(pokemon, p.Pokemon.Name)
	}

	return pokemon, nil
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
