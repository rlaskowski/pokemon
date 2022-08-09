package pokemon

import "github.com/labstack/echo/v4/middleware"

func CorsConfig() middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}
}
