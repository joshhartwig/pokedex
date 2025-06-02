package main

import (
	"time"

	"github.com/joshhartwig/pokedex/internal/pokecache"
	"github.com/joshhartwig/pokedex/internal/repl"
	"github.com/joshhartwig/pokedex/pkg/models"
)

func main() {
	app := models.Config{}
	app.Commands = map[string]models.CliCommand{
		"exit": {
			Name:        "exit",
			Description: "used to close the app",
			Callback:    func(args ...string) error { return repl.ExitCmd(&app, args...) },
		},
		"help": {
			Name:        "help",
			Description: "used to get help",
			Callback:    func(args ...string) error { return repl.HelpCmd(&app, args...) },
		},
		"map": {
			Name:        "map",
			Description: "used to list all the pokedex locations",
			Callback:    func(args ...string) error { return repl.MapCmd(&app, args...) },
		},
		"mapb": {
			Name:        "mapb",
			Description: "used to move forward in the map",
			Callback:    func(args ...string) error { return repl.MapbCmd(&app, args...) },
		},
		"explore": {
			Name:        "explore",
			Description: "explores a section of the map",
			Callback:    func(args ...string) error { return repl.AltExploreCmd(&app, args...) },
		},
		"catch": {
			Name:        "catch",
			Description: "attempts to catch a pokemon",
			Callback:    func(args ...string) error { return repl.Catch(&app, args...) },
		},
	}
	app.BaseApiUrl = "https://pokeapi.co/api/v2/location-area/"

	// create a new cache with a timer of 10 seconds
	app.Cache = *pokecache.NewCache(time.Millisecond * 10)

	// start the repl loop
	repl.Repl(&app)
}
