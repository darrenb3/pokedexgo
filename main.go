package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// A Response struct to map the Pokemon's data to
type Pokemon struct {
	Name    string  `json:"name"`
	ID      int     `json:"id"`
	Weight  float64 `json:"weight"`
	Sprites Sprites `json:"sprites"`
	Stats   []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
}

// Struct to map sprites of the pokemon
type Sprites struct {
	FrontDefault string `json:"front_default"`
}

// Style definitons
var textStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#7D56F4"))

var tableTextStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#ffa6c9")).PaddingLeft(1).PaddingRight(1)

var warnStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#CD5C5C"))

func main() {
main:
	for {
		var userInput string
		fmt.Println(textStyle.Render("Please input a Pokémon name or id: "))
		fmt.Scanln(&userInput)
		switch userInput {
		case "exit":
			fmt.Print("Exiting PokédexGO...")
			break main
		case "":
			fmt.Println(warnStyle.Render("Please enter the name or ID of a Pokémon!"))
		default:
			url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", userInput)
			response, err := http.Get(url)

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
			if responseObject.ID == 0 {
				fmt.Println(warnStyle.Render("Pokémon not found. Please try again..."))
			} else {
				var stats []int
				for _, stat := range responseObject.Stats {
					stats = append(stats, stat.BaseStat)
				}
				rows := [][]string{
					{"Name:", responseObject.Name},
					{"ID:", fmt.Sprintf("%d", responseObject.ID)},
					{"HP:", fmt.Sprintf("%d", stats[0])},
					{"Attack:", fmt.Sprintf("%d", stats[1])},
					{"Defense:", fmt.Sprintf("%d", stats[2])},
					{"Sp. Atk:", fmt.Sprintf("%d", stats[3])},
					{"Sp. Def:", fmt.Sprintf("%d", stats[4])},
					{"Speed:", fmt.Sprintf("%d", stats[5])},
					//
				}

				t := table.New().
					Border(lipgloss.NormalBorder()).
					BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
					StyleFunc(func(row, col int) lipgloss.Style {
						return tableTextStyle
					}).
					Rows(rows...)
				fmt.Println(t)
			}
		}
	}
}
