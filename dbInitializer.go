package dbInitializer

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

type allPokemon struct {
	Count    int `json:"count"`
	Next     any `json:"next"`
	Previous any `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// A Response struct to map the Pokemon's data to
type Pokemon struct {
	Name   string  `json:"name"`
	ID     int     `json:"id"`
	Weight float64 `json:"weight"`
	Stats  []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
	}
}

func getAllPokemon() allPokemon {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon?limit=100000&offset=0")
	response, err := http.Get(url)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject allPokemon
	json.Unmarshal(responseData, &responseObject)

	return responseObject
}

func getPokemonUrl(allPokemon allPokemon) []string {
	var urls []string
	for _, result := range allPokemon.Results {
		urls = append(urls, result.URL)
	}
	return urls
}
