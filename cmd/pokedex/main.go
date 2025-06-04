package main

import (
	"database/sql"
	"fmt"
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
[ ] - setup standard logger and clean up logging
[x] - setup database and goose + sqlc
[ ] - setup a way to delete pokemon
[ ] - change catch algo to factor in skill
[ ] -
 - change the add functionality to add pokemon to database
- fetch json data from db
- get pokemon to fight via fight command
- change catch diff to factor in skill

*/

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("unable to load environment variables... exiting")
		os.Exit(1) // quit to os
	}

	dbConnStr := os.Getenv("POSTGRES_CONNSTR")
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		fmt.Println("error connecting to db")
		fmt.Println(err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	app := models.Config{}
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
			Callback:    func(args ...string) error { return repl.AltExplore(&app, args...) },
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
	}
	app.Pokedex = map[string]models.Pokemon{}
	app.BaseApiUrl = "https://pokeapi.co/api/v2/location-area/"
	app.Db = dbQueries

	// create a new cache with a timer of 10 seconds
	app.Cache = *pokecache.NewCache(time.Millisecond * 10)

	// start the repl loop
	repl.Repl(&app)
}
