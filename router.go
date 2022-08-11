package pokemon

import (
	"log"
	"net/http"
	"reflect"
	"rlaskowski/pokemon/cmd"
	"strings"

	"github.com/labstack/echo/v4"
)

type Router struct {
	echo    *echo.Echo
	pokemon *Pokemon
}

func NewRouter(echo *echo.Echo) *Router {
	return &Router{
		echo:    echo,
		pokemon: NewPokemon(cmd.ServiceAddr),
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
	pr := &PokemonRequest{}

	if err := ctx.Bind(pr); err != nil || reflect.DeepEqual(pr, &PokemonRequest{}) {
		return r.errorHandler(ErrParameter, ctx)
	}

	result, err := r.pokemon.Fight(pr)
	if err != nil {
		return r.errorHandler(ErrContent, ctx)
	}

	return ctx.JSON(http.StatusOK, result)
}

func (r *Router) Pokemon(ctx echo.Context) error {
	name := ctx.QueryParam("pokemonName")
	if len(name) > 0 {
		pi, err := r.pokemon.PokemonType(name)

		if err != nil {
			log.Printf("pokemon type error: %s", err.Error())

			return r.errorHandler(ErrContent, ctx)
		}

		return ctx.JSON(http.StatusOK, pi)
	}

	pn, err := r.pokemon.Name()
	if err != nil {
		log.Printf("pokemon name error: %s", err.Error())

		return r.errorHandler(ErrContent, ctx)
	}

	return ctx.JSON(http.StatusOK, pn)
}

func (r *Router) Weak(ctx echo.Context) error {
	return r.weak(true, ctx)
}

func (r *Router) Strong(ctx echo.Context) error {
	return r.weak(false, ctx)
}

func (r *Router) weak(weak bool, ctx echo.Context) error {
	name := strings.TrimSpace(ctx.QueryParam("typeName"))
	if !(len(name) > 0) {
		return r.errorHandler(ErrParameter, ctx)
	}

	list, err := r.pokemon.WeakOrStrong(name, weak)
	if err != nil {
		return r.errorHandler(ErrContent, ctx)
	}

	return ctx.JSON(http.StatusOK, list)
}

func (r *Router) errorHandler(err Error, ctx echo.Context) error {
	return ctx.JSON(err.Code, err)
}
