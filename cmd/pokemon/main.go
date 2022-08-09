package main

import (
	"log"
	"os"
	"rlaskowski/pokemon"
	"rlaskowski/pokemon/cmd"
)

func main() {

	cmd.RunFlags()
	service := pokemon.NewService()

	if err := service.Run(); err != nil {
		log.Fatalf("Unexpected error: %s", err.Error())
		os.Exit(1)
	}
}
