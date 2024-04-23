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
main:
	for {
		var userInput string
		fmt.Print("Please input a Pokémon name or id: ")
		fmt.Scanln(&userInput)
		switch userInput {
		case "exit":
			fmt.Print("Exiting PokédexGO...")
			break main
		case "":
			fmt.Print("Please enter the name or ID of a Pokémon!\n")
		default:
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
			if responseObject.ID == 0 {
				fmt.Print("Pokémon not found. Please try again...\n")
			} else {
				fmt.Printf("Pokémon name is: %s\n", responseObject.Name)
				fmt.Printf("Pokémon id is: %d\n", responseObject.ID)
				fmt.Printf("Pokémon weighs: %.1f kg\n", responseObject.Weight*.1)
				fmt.Printf("Link to Pokémon's sprite:\n%s", responseObject.Sprites.FrontDefault)
			}
		}
	}
}
