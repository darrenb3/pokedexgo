package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

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

// makes the first letter of an ascii string uppercase
func upperFirstLetter(s string) string {
	letter := string(s[0])
	uppperLetter := strings.ToUpper(letter)
	newString := strings.Replace(s, letter, uppperLetter, 1)
	return newString
}

func main() {
main:
	for {
		var userInput string
		fmt.Println(textStyle.Render("Please input a Pokémon name or id: "))
		fmt.Scanln(&userInput)
		userInput = strings.ToLower(userInput)
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

				rows := [][]string{
					{"Name:", upperFirstLetter(responseObject.Name)},
					{"ID:", fmt.Sprintf("%d", responseObject.ID)},
					{"Type(s):", finalType},
					{"HP:", fmt.Sprintf("%d", stats[0])},
					{"Attack:", fmt.Sprintf("%d", stats[1])},
					{"Defense:", fmt.Sprintf("%d", stats[2])},
					{"Sp. Atk:", fmt.Sprintf("%d", stats[3])},
					{"Sp. Def:", fmt.Sprintf("%d", stats[4])},
					{"Speed:", fmt.Sprintf("%d", stats[5])},
				}

				t := table.New().
					Border(lipgloss.NormalBorder()).
					BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
					StyleFunc(func(row, col int) lipgloss.Style {
						return tableTextStyle
					}).
					Rows(rows...)

				//Making clickable link to sprite
				spriteLink := textStyle.Render(fmt.Sprintf("\x1b]8;;%s\x07Link to Pokemon's sprite\x1b]8;;\x07\u001b[0m", responseObject.Sprites.FrontDefault))

				fmt.Println(t)
				fmt.Println(spriteLink)
			}
		}
	}
}
