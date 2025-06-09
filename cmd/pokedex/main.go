package main

import (
	"context"
	"database/sql"
	"encoding/json"
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
		"fight": {
			Name:        "fight",
			Description: "fight two pokemon that you have caught",
			Callback:    func(args ...string) error { return repl.Fight(&conf, args...) },
		},
	}
	loadPokedexFromDb(&conf)
	repl.Repl(&conf)
}

// loadPokedexFromDb loads Pokemon data from the database into the Pokedex cache.
// It retrieves all Pokemon rows from the database, unmarshals the JSON data for each Pokemon,
// and stores them in the Config's Pokedex map using the Pokemon name as the key.
// If there is an error retrieving the Pokemon list or unmarshaling individual Pokemon data,
// it will either return the error or log it and continue to the next Pokemon.
//
// Parameters:
//   - c: A pointer to the Config struct containing the database connection and Pokedex map
//
// Returns:
//   - error: Returns an error if the database query fails, nil otherwise
func loadPokedexFromDb(c *models.Config) error {
	pokemonRows, err := c.Db.ListPokemon(context.Background())
	if err != nil {
		return err
	}

	for _, row := range pokemonRows {
		var p models.Pokemon
		if err := json.Unmarshal(row.JsonData.RawMessage, &p); err != nil {
			c.Logger.Error("failed to unmarhsal pokemon")
			continue
		}

		c.Pokedex[row.PokemonName] = p

	}
	return nil
}
