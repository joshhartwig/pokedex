package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/joshhartwig/pokedex/internal/database"
	"github.com/joshhartwig/pokedex/internal/pokecache"
	"github.com/joshhartwig/pokedex/internal/repl"
	"github.com/joshhartwig/pokedex/pkg/models"

	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("unable to load environment variables: %v\n", err)
		os.Exit(1) // quit to os
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	dbConnStr := os.Getenv("POSTGRES_CONNSTR")
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		fmt.Printf("error connecting to db: %v\n", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	conf := models.Config{
		Logger:     logger,
		Db:         dbQueries,
		BaseApiUrl: "https://pokeapi.co/api/v2/location-area/",
		Cache:      *pokecache.NewCache(time.Millisecond * 10),
		Pokedex:    map[string]models.Pokemon{},
	}

	// setup the commands
	conf.Commands = map[string]models.CliCommand{
		"exit": {
			Name:        "exit",
			Description: "used to close the app",
			Callback:    func(args ...string) error { return repl.Exit(&conf, args...) },
		},
		"help": {
			Name:        "help",
			Description: "used to get help",
			Callback:    func(args ...string) error { return repl.Help(&conf, args...) },
		},
		"map": {
			Name:        "map",
			Description: "used to list all the pokedex locations",
			Callback:    func(args ...string) error { return repl.Map(&conf, args...) },
		},
		"mapb": {
			Name:        "mapb",
			Description: "used to move forward in the map",
			Callback:    func(args ...string) error { return repl.Mapb(&conf, args...) },
		},
		"explore": {
			Name:        "explore",
			Description: "explores a section of the map",
			Callback:    func(args ...string) error { return repl.Explore(&conf, args...) },
		},
		"catch": {
			Name:        "catch",
			Description: "attempts to catch a pokemon",
			Callback:    func(args ...string) error { return repl.Catch(&conf, args...) },
		},
		"inspect": {
			Name:        "inspect",
			Description: "inspects a caught pokemon",
			Callback:    func(args ...string) error { return repl.Inspect(&conf, args...) },
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "displays your caught pokemons",
			Callback:    func(args ...string) error { return repl.Pokedex(&conf, args...) },
		},
		"history": {
			Name:        "history",
			Description: "display command line history for each command",
			Callback:    func(args ...string) error { return repl.History(&conf, args...) },
		},
	}

	repl.Repl(&conf)
}
