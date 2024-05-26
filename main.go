package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"github.com/cavaliergopher/grab/v3"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/qeesung/image2ascii/convert"
	_ "modernc.org/sqlite"
)

const database string = "pokemon.db"

type pokemonStruct struct {
	Id         int
	Name       string
	Types      string
	Hp         int
	Attack     int
	Defense    int
	Sp_atk     int
	Sp_def     int
	Speed      int
	Sprite_URL string
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

// Downloads the image from the api url linked and then converts it to ascii art for display in tui
func imageArt(spriteURL string) string {

	//downloading image from url
	resp, err := grab.Get(".", spriteURL)
	if err != nil {
		log.Fatal(err)
	}
	convertOptions := convert.DefaultOptions
	convertOptions.FitScreen = true
	converter := convert.NewImageConverter()
	image := converter.ImageFile2ASCIIString(resp.Filename, &convertOptions)
	os.Remove(resp.Filename)

	//Cleaning the ascii art image to remove the png transparency padding that is added on conversion
	var newImage []string
	temp := strings.Split(image, "\n")
	for _, lineContent := range temp {
		line := strings.ReplaceAll(lineContent, "\x1b[0;00m", " ")
		line = strings.ReplaceAll(line, "\x1b[38;5;16m", " ")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		} else {
			newImage = append(newImage, lineContent)
		}
	}
	return strings.Join(newImage, "\n")
}

func checkDbExistence() error {
	if _, err := os.Stat("pokemon.db"); err == nil {
		log.Printf("Database already exist\n")
	} else {
		log.Printf("Database not created...\n")
		log.Printf("Creating Database...\n")
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
	}
	return nil
}

func main() {

	checkDbExistence()

	//Database connection
	db, err := sql.Open("sqlite", database)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = createDatabase(db); err != nil {
		fmt.Println(err)
		return
	}

	should_exit := false

	for {
		if should_exit {
			break
		}
		var userInput string
		fmt.Println(textStyle.Render("Please input a Pokémon name or id: "))
		fmt.Scanln(&userInput)
		userInput = strings.ToLower(userInput)
		switch userInput {
		case "exit":
			fmt.Print("Exiting PokédexGO...")
			should_exit = true
		case "":
			fmt.Println(warnStyle.Render("Please enter the name or ID of a Pokémon!"))
		default:
			var pokemon pokemonStruct
			query := `SELECT * FROM pokemon WHERE name LIKE ? limit 1`
			row := db.QueryRow(query, userInput)
			if err = row.Scan(&pokemon.Id, &pokemon.Name, &pokemon.Types, &pokemon.Hp, &pokemon.Attack, &pokemon.Defense, &pokemon.Sp_atk, &pokemon.Sp_def, &pokemon.Speed, &pokemon.Sprite_URL); err == sql.ErrNoRows {
				log.Printf("Pokemon not found")
				continue
			} else if err != nil {
				log.Println(err)
				continue
			}
			rows := [][]string{
				{"Name:", upperFirstLetter(pokemon.Name)},
				{"ID:", fmt.Sprintf("%d", pokemon.Id)},
				{"Type(s):", pokemon.Types},
				{"HP:", fmt.Sprintf("%d", pokemon.Hp)},
				{"Attack:", fmt.Sprintf("%d", pokemon.Attack)},
				{"Defense:", fmt.Sprintf("%d", pokemon.Defense)},
				{"Sp. Atk:", fmt.Sprintf("%d", pokemon.Sp_atk)},
				{"Sp. Def:", fmt.Sprintf("%d", pokemon.Sp_def)},
				{"Speed:", fmt.Sprintf("%d", pokemon.Speed)},
			}

			t := table.New().
				Border(lipgloss.NormalBorder()).
				BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
				StyleFunc(func(row, col int) lipgloss.Style {
					return tableTextStyle
				}).
				Rows(rows...)

			//Making clickable link to sprite
			spriteLink := textStyle.Render(fmt.Sprintf("\x1b]8;;https://www.serebii.net/pokemon/%s\x07Link to Pokemon's Serebii.net entry\x1b]8;;\x07\u001b[0m", pokemon.Name))

			//Printing Pokemon's info
			fmt.Println(imageArt(pokemon.Sprite_URL))
			fmt.Println(t)
			fmt.Println(spriteLink)

		}
	}
}
