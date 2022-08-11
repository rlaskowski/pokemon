package pokemon

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	WinResult  = "win"
	LoseResult = "lose"
	DrawResult = "draw"
)

type PokemonsService struct {
	Pokemon []PokemonService `json:"pokemon"`
}

type PokemonName struct {
	Name string `json:"name"`
}

type PokemonType struct {
	Image       string   `json:"image"`
	TypeName    []string `json:"type"`
	SpawnChance float64  `json:"spawn_chance"`
	Weaknesses  []string `json:"-"`
	Name        string   `json:"-"`
}

type NextEvolution struct {
	Num  string `json:"num"`
	Name string `json:"name"`
}

type PokemonService struct {
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

type PokemonRequest struct {
	Owner string `json:"myPokemon" form:"myPokemon"`
	Enemy string `json:"enemyPokemon" form:"enemyPokemon"`
}

type Pokemon struct {
	serviceURL string
	pool       sync.Pool
}

func NewPokemon(url string) *Pokemon {
	p := &Pokemon{
		serviceURL: url,
	}

	p.pool.New = func() interface{} {
		return &PokemonsService{}
	}

	return p
}

func (p *Pokemon) Fight(preq *PokemonRequest) (string, error) {
	pt, err := p.PokemonType(preq.Owner, preq.Enemy)
	if err != nil {
		log.Printf("pokemons error: %s", err.Error())

		return "", err
	}

	if !(len(pt) == 2) {
		err = errors.New("incorrect pokemon list")
		log.Printf("pokemons error: %s", err.Error())

		return "", err
	}

	return p.fight(pt)
}

func (p *Pokemon) fight(list []PokemonType) (string, error) {
	pm := p.fightResult(list[0].TypeName, list[1].Weaknesses)
	pe := p.fightResult(list[1].TypeName, list[0].Weaknesses)

	if strings.Compare(pm, WinResult) == 0 {
		return pm, nil
	}

	if strings.Compare(pe, WinResult) == 0 {
		return LoseResult, nil
	}

	return DrawResult, nil
}

func (p *Pokemon) fightResult(names []string, weaknesses []string) string {
	for _, n := range names {
		if eq := p.findName(n, weaknesses); eq {
			return "win"
		}
	}

	return "draw"
}

func (p *Pokemon) Name() (list []PokemonName, err error) {
	psrv, err := p.pokemonService()

	if err != nil {
		log.Printf("pokemons error: %s", err.Error())

		return nil, err
	}

	for _, p := range psrv.Pokemon {
		pn := PokemonName{
			Name: p.Name,
		}

		list = append(list, pn)
	}

	return list, nil
}

func (p *Pokemon) PokemonType(names ...string) (list []PokemonType, err error) {
	psrv, err := p.pokemonService()

	if err != nil {
		log.Printf("pokemons error: %s", err.Error())

		return nil, err
	}

	for _, name := range names {
		for _, p := range psrv.Pokemon {
			name = strings.TrimSpace(name)
			p.Name = strings.TrimSpace(p.Name)

			if eq := bytes.EqualFold([]byte(name), []byte(p.Name)); eq {
				pi := PokemonType{
					Image:       p.Image,
					TypeName:    p.TypeName,
					SpawnChance: p.SpawnChance,
					Weaknesses:  p.Weaknesses,
					Name:        p.Name,
				}

				list = append(list, pi)
			}
		}
	}

	return list, nil
}

func (p *Pokemon) WeakOrStrong(name string, weak bool) (list PokemonsService, err error) {
	psrv, err := p.pokemonService()

	if err != nil {
		log.Printf("pokemons error: %s", err.Error())

		return PokemonsService{}, err
	}

	for _, ps := range psrv.Pokemon {
		eq := p.findName(name, ps.Weaknesses)

		if !(eq == weak) {
			continue
		}

		list.Pokemon = append(list.Pokemon, ps)
	}

	return list, nil
}

func (p *Pokemon) findName(name string, list []string) bool {
	name = strings.TrimSpace(name)

	for _, l := range list {
		l = strings.TrimSpace(l)

		if eq := bytes.EqualFold([]byte(l), []byte(name)); eq {
			return true
		}
	}

	return false
}

func (p *Pokemon) pokemonService() (*PokemonsService, error) {
	psrv := p.pool.Get().(*PokemonsService)

	defer p.pool.Put(psrv)

	res, err := http.Get(p.serviceURL)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, psrv); err != nil {
		return nil, err
	}

	return psrv, nil
}
