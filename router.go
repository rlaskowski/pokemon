package pokemon

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"rlaskowski/pokemon/cmd"
	"strings"

	"github.com/labstack/echo/v4"
)

type Router struct {
	echo    *echo.Echo
	pokemon Pokemons
}

func NewRouter(echo *echo.Echo) *Router {
	return &Router{
		echo: echo,
	}
}

// Initializing all endpoints
func (r *Router) Run() {
	r.echo.GET("/", r.Pokemon)
	r.echo.GET("/weak/", r.Weak)
	r.echo.GET("/strong/", r.Strong)
	r.echo.POST("/figth/", r.Fight)
}

func (r *Router) Fight(ctx echo.Context) error {
	return nil
}

func (r *Router) Pokemon(ctx echo.Context) error {
	if err := r.pokemons(); err != nil {
		log.Printf("pokemon name error: %s", err.Error())

		return r.errorHandler(ErrContent, ctx)
	}

	name := ctx.QueryParam("pokemonName")
	if len(name) > 0 {
		pi, err := r.image(name)

		if err != nil {
			log.Printf("pokemon name error: %s", err.Error())

			return r.errorHandler(ErrContent, ctx)
		}

		return ctx.JSON(http.StatusOK, pi)
	}

	pn, err := r.name()
	if err != nil {
		log.Printf("pokemon name error: %s", err.Error())

		return r.errorHandler(ErrContent, ctx)
	}

	return ctx.JSON(http.StatusOK, pn)
}

func (r *Router) Weak(ctx echo.Context) error {
	return r.weakOrStrong(ctx, true)
}

func (r *Router) Strong(ctx echo.Context) error {
	return r.weakOrStrong(ctx, false)
}

// Refreshing pokemon list
func (r *Router) pokemons() error {
	res, err := http.Get(cmd.ServiceAddr)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	ps := &Pokemons{}

	if err := json.Unmarshal(body, ps); err != nil {
		return err
	}

	r.pokemon = *ps

	return nil

}

func (r *Router) errorHandler(err Error, ctx echo.Context) error {
	return ctx.JSON(err.Code, err)
}

func (r *Router) name() ([]PokemonName, error) {
	pc := make([]PokemonName, 0)

	for _, p := range r.pokemon.Pokemon {
		pn := PokemonName{
			Name: p.Name,
		}

		pc = append(pc, pn)
	}

	return pc, nil
}

func (r *Router) image(name string) ([]PokemonImage, error) {
	pc := make([]PokemonImage, 0)

	for _, p := range r.pokemon.Pokemon {
		name = strings.TrimSpace(strings.ToLower(name))
		p.Name = strings.TrimSpace(strings.ToLower(p.Name))

		if strings.Compare(p.Name, name) == 0 {
			pi := PokemonImage{
				Image:       p.Image,
				TypeName:    p.TypeName,
				SpawnChance: p.SpawnChance,
			}

			pc = append(pc, pi)
		}
	}

	return pc, nil
}

// Returns pokemons that are weak or strong - according to weak parameter
func (r *Router) weakOrStrong(ctx echo.Context, weak bool) error {
	name := ctx.QueryParam("typeName")
	if !(len(name) > 0) {
		return r.errorHandler(ErrQueryParameter, ctx)
	}

	if err := r.pokemons(); err != nil {
		log.Printf("pokemon name error: %s", err.Error())

		return r.errorHandler(ErrContent, ctx)
	}

	pc := Pokemons{}

	for _, p := range r.pokemon.Pokemon {
		for _, w := range p.Weaknesses {
			w = strings.TrimSpace(strings.ToLower(w))
			name = strings.TrimSpace(strings.ToLower(name))

			weakC := strings.Compare(w, name)

			if (weakC == 0 && weak) || (weakC != 0 && !weak) {
				pc.Pokemon = append(pc.Pokemon, p)
			}
		}
	}

	return ctx.JSON(http.StatusOK, pc)
}
