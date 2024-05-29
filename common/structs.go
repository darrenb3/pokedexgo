package common

// Struct for holding the information of a Pokemon
type Pokemon struct {
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

// Struct for holding response of all pokemon query
type AllPokemon struct {
	Count    int `json:"count"`
	Next     any `json:"next"`
	Previous any `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// A Response struct to map the Pokemon's data to from the api call
type PokemonReponse struct {
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
