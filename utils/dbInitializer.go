package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"pokedexgo/common"
	"strings"

	_ "modernc.org/sqlite"
)

// makes the first letter of an ascii string uppercase
func upperFirstLetter(s string) string {
	letter := string(s[0])
	uppperLetter := strings.ToUpper(letter)
	newString := strings.Replace(s, letter, uppperLetter, 1)
	return newString
}

func createDatabase(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS pokemon (id INT PRIMARY KEY, name STRING, type STRING, hp INT, attack INT, defense INT, sp_atk INT, sp_def INT, speed INT, sprite_url STRING)`
	if _, err := db.Exec(query); err != nil {
		fmt.Println("Error in table create")
		return err
	}
	return nil
}

func getAllPokemon() common.AllPokemon {
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

	var responseObject common.AllPokemon
	json.Unmarshal(responseData, &responseObject)

	return responseObject
}

func getPokemonUrl(allPokemon common.AllPokemon) []string {
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
		var responseObject common.PokemonReponse
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
		pokemon := common.Pokemon{Id: responseObject.ID, Name: upperFirstLetter(responseObject.Name), Types: finalType, Hp: stats[0], Attack: stats[1], Defense: stats[2], Sp_atk: stats[3], Sp_def: stats[4], Speed: stats[5], Sprite_URL: responseObject.Sprites.FrontDefault}
		query := `INSERT INTO POKEMON (id, name, type, hp, attack, defense, sp_atk, sp_def, speed, sprite_url) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		if _, err := db.Exec(query, pokemon.Id, pokemon.Name, pokemon.Types, pokemon.Hp, pokemon.Attack, pokemon.Defense, pokemon.Sp_atk, pokemon.Sp_def, pokemon.Speed, pokemon.Sprite_URL); err != nil {
			fmt.Println("Error in pokemon insert")
			return err
		}

	}
	return nil
}

func DbInitialize() error {
	var database string = "pokemon.db"
	db, err := sql.Open("sqlite", database)
	if err != nil {
		fmt.Println(err)
	}
	if err = createDatabase(db); err != nil {
		fmt.Println(err)
	}
	createDatabase(db)
	pokemons := getAllPokemon()
	pokemonURLS := getPokemonUrl(pokemons)
	addPokemonToDatabase(db, pokemonURLS)
	return nil
}
