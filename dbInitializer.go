package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

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

// makes the first letter of an ascii string uppercase
func upperFirstLetter(s string) string {
	letter := string(s[0])
	uppperLetter := strings.ToUpper(letter)
	newString := strings.Replace(s, letter, uppperLetter, 1)
	return newString
}

func createDatabase(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS pokemon (id INT PRIMARY KEY, name STRING, type STRING, hp INT, attack INT, defense INT, sp_atk INT, sp_def INT, speed INT)`
	if _, err := db.Exec(query); err != nil {
		fmt.Println("Error in table create")
		return err
	}
	return nil
}

func getAllPokemon() allPokemon {
	url := "https://pokeapi.co/api/v2/pokemon?limit=100000&offset=0"
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

func addPokemonToDatabase(db *sql.DB, urls []string) error {
	for _, pokemonURL := range urls {
		response, err := http.Get(pokemonURL)

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		var responseObject Pokemon
		json.Unmarshal(responseData, &responseObject)
		var stats []int
		var types []string
		//Mapping all stats to a slice
		for _, stat := range responseObject.Stats {
			stats = append(stats, stat.BaseStat)
		}
		//Mapping all types to a slice
		for _, pTypes := range responseObject.Types {
			types = append(types, pTypes.Type.Name)
		}
		//Checking if pokemon has 2 types or not: If yes combines them into a single string
		var finalType string
		if len(types) == 2 {
			finalType = strings.ToUpper(strings.Join(types, " "))
		} else {
			finalType = strings.ToUpper(types[0])
		}
		var pokemonData []string
		pokemonData = append(pokemonData, fmt.Sprintf("%d", responseObject.ID), upperFirstLetter(responseObject.Name), finalType, fmt.Sprintf("%d", stats[0]), fmt.Sprintf("%d", stats[1]), fmt.Sprintf("%d", stats[2]), fmt.Sprintf("%d", stats[3]), fmt.Sprintf("%d", stats[4]), fmt.Sprintf("%d", stats[5]))
		query := `INSERT INTO POKEMON (id, name, type, hp, attack, defense, sp_atk, sp_def, speed) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`
		if _, err := db.Exec(query, pokemonData[0], pokemonData[1], pokemonData[2], pokemonData[3], pokemonData[4], pokemonData[5], pokemonData[6], pokemonData[7], pokemonData[8]); err != nil {
			fmt.Println("Error in pokemon insert")
			return err
		}

	}
	return nil
}
