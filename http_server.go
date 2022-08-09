package pokemon

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"rlaskowski/pokemon/cmd"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HttpServer struct {
	echo    *echo.Echo
	context context.Context
	router  *Router
}

func NewHttpServer(ctx context.Context) *HttpServer {

	echo := echo.New()

	return &HttpServer{
		echo:    echo,
		context: ctx,
		router:  NewRouter(echo),
	}
}

// Starting http server
func (h *HttpServer) Start() error {
	h.echo.HideBanner = true
	h.echo.HidePort = true
	h.echo.Use(middleware.CORSWithConfig(CorsConfig()))
	h.echo.Use(middleware.Recover())

	h.router.Run()

	log.Printf("Starting HTTP server on http://localhost:%d", cmd.HttpPort)

	if err := h.echo.Start(fmt.Sprintf(":%d", cmd.HttpPort)); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("caught error while starting server: %s", err.Error())
	}

	return nil
}

// Stopping http server
func (h *HttpServer) Stop() error {
	log.Print("Stopping HTTP server")

	return h.echo.Close()
}
