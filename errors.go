package pokemon

import "net/http"

var (
	ErrQueryParameter = newError(http.StatusBadRequest, "bad query parameter")
	ErrContent        = newError(http.StatusBadRequest, "content wasn't found")
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
