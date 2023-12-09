package pokeshell

import (
	"io"
	"net/http"
)

type LocationPage struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type PokemonEncounter struct {
	Pokemon struct {
		Name string `json:"name"`
	} `json:"pokemon"`
}

type PokemonFetched struct {
	Name   string  `json:"name"`
	Height float64 `json:"height"`
	Weight float64 `json:"weight"`
	Types  []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
}

type LocationAreaExplore struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

func fetchAPI(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
