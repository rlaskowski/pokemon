package pokemon

import "net/http"

var (
	ErrParameter   = newError(http.StatusBadRequest, "bad request parameter")
	ErrContent     = newError(http.StatusBadRequest, "content wasn't found")
	ErrPokemonName = newError(http.StatusBadRequest, "bad pokemon name list")
)

type Error struct {
	Name string `json:"error_name"`
	Code int    `json:"code"`
}

func newError(code int, name string) Error {
	return Error{
		Name: name,
		Code: code,
	}
}
