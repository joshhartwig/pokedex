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

/*
TODO:
[ ] - add check to catch so we don't catch one with the same name and error testing this logic, it should fail as we have a unique name constraint
[x] - setup standard logger and clean up logging
[x] - setup database and goose + sqlc
[ ] - setup a way to delete pokemon
[ ] - change catch algo to factor in skill
[ ] - add a fight command to fight pokemon
[ ] - fix the bug where if you call mapb before map it will not work

*/

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

	app := models.Config{
		Logger:     logger,
		Db:         dbQueries,
		BaseApiUrl: "https://pokeapi.co/api/v2/location-area/",
		Cache:      *pokecache.NewCache(time.Millisecond * 10),
		Pokedex:    map[string]models.Pokemon{},
	}

	app.Commands = map[string]models.CliCommand{
		"exit": {
			Name:        "exit",
			Description: "used to close the app",
			Callback:    func(args ...string) error { return repl.Exit(&app, args...) },
		},
		"help": {
			Name:        "help",
			Description: "used to get help",
			Callback:    func(args ...string) error { return repl.Help(&app, args...) },
		},
		"map": {
			Name:        "map",
			Description: "used to list all the pokedex locations",
			Callback:    func(args ...string) error { return repl.Map(&app, args...) },
		},
		"mapb": {
			Name:        "mapb",
			Description: "used to move forward in the map",
			Callback:    func(args ...string) error { return repl.Mapb(&app, args...) },
		},
		"explore": {
			Name:        "explore",
			Description: "explores a section of the map",
			Callback:    func(args ...string) error { return repl.Explore(&app, args...) },
		},
		"catch": {
			Name:        "catch",
			Description: "attempts to catch a pokemon",
			Callback:    func(args ...string) error { return repl.Catch(&app, args...) },
		},
		"inspect": {
			Name:        "inspect",
			Description: "inspects a caught pokemon",
			Callback:    func(args ...string) error { return repl.Inspect(&app, args...) },
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "displays your caught pokemons",
			Callback:    func(args ...string) error { return repl.Pokedex(&app, args...) },
		},
		"history": {
			Name:        "history",
			Description: "display command line history for each command",
			Callback:    func(args ...string) error { return repl.History(&app, args...) },
		},
	}

	repl.Repl(&app)
}
