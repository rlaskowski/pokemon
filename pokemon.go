package pokemon

type Pokemons struct {
	Pokemon []Pokemon `json:"pokemon"`
}

type PokemonName struct {
	Name string `json:"name"`
}

type PokemonImage struct {
	Image       string   `json:"image"`
	TypeName    []string `json:"type"`
	SpawnChance float64  `json:"spawn_chance"`
}

type NextEvolution struct {
	Num  string `json:"num"`
	Name string `json:"name"`
}

type Pokemon struct {
	ID            int             `json:"id"`
	Number        string          `json:"num"`
	Name          string          `json:"name" query:"pokemonName"`
	Image         string          `json:"img"`
	TypeName      []string        `json:"type" query:"typeName"`
	Height        string          `json:"height"`
	Weight        string          `json:"weight"`
	Candy         string          `json:"candy"`
	CandyCount    int             `json:"candy_count"`
	Egg           string          `json:"egg"`
	SpawnChance   float64         `json:"spawn_chance"`
	AvgSpawns     float64         `json:"avg_spawns"`
	SpawnTime     string          `json:"spawn_time"`
	Multipliers   []float64       `json:"multipliers"`
	Weaknesses    []string        `json:"weaknesses"`
	NextEvolution []NextEvolution `json:"next_evolution"`
}
