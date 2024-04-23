package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// A Response struct to map the Pokemon's data to
type Pokemon struct {
	Name    string  `json:"name"`
	ID      int     `json:"id"`
	Weight  float64 `json:"weight"`
	Sprites Sprites `json:"sprites"`
}

// Struct to map sprites of the pokemon
type Sprites struct {
	FrontDefault string `json:"front_default"`
}

func main() {
	var userInput string
	fmt.Print("Please input a Pokémon name or id: ")
	fmt.Scan(&userInput)
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", userInput)
	response, err := http.Get(url)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseObject Pokemon
	json.Unmarshal(responseData, &responseObject)

	fmt.Printf("Pokémon name is: %s\n", responseObject.Name)
	fmt.Printf("Pokémon id is: %d\n", responseObject.ID)
	fmt.Printf("Pokémon weighs: %.1f kg\n", responseObject.Weight*.1)
	fmt.Printf("Link to Pokémon's sprite:\n%s", responseObject.Sprites.FrontDefault)

}
