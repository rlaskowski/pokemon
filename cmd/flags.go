package cmd

import "flag"

var (
	HttpPort    int
	ServiceAddr string
)

func RunFlags() {
	flag.IntVar(&HttpPort, "p", 8080, "port for http server")
	flag.StringVar(&ServiceAddr, "a", "https://raw.githubusercontent.com/Biuni/PokemonGO-Pokedex/master/pokedex.json", "service address")

	flag.Parse()
}
