package repl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/joshhartwig/pokedex/internal/api"
	"github.com/joshhartwig/pokedex/internal/database"
	"github.com/joshhartwig/pokedex/pkg/models"
	"github.com/sqlc-dev/pqtype"
)

// Exit is a command that exits the Pokedex application.
func Exit(c *models.Config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// Help displays the help screen for the Pokedex application.
func Help(c *models.Config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Printf("  %-10s %s\n", "help:", "Displays this help screen")
	fmt.Printf("  %-10s %s\n", "map:", "Displays map of locations from api, consecutive calls move the map forward")
	fmt.Printf("  %-10s %s\n", "mapb:", "Displays previous map, if called with no map call, it will call the base api")
	fmt.Printf("  %-10s %s\n", "catch", "Attempts to catch a specific pokemon")
	fmt.Printf("  %-10s %s\n", "explore:", "Displays pokemon in a specific region")
	fmt.Printf("  %-10s %s\n", "inspect:", "Displays stats for a specific pokemon (must be caught first)")
	fmt.Printf("  %-10s %s\n", "exit:", "Exit the Pokedex")
	return nil
}

// Map fetches the map of locations from the PokeAPI and displays them.
func Map(c *models.Config, args ...string) error {

	var ah models.Apiheader
	// if c.next is anything but empty, it likely has a url and pull from that url
	if c.Next != "" {
		if err := api.FetchFromCache(c, c.Next, &ah); err != nil {
			c.Logger.Error("error fetching from cache", "error", err)
			return err
		}
	} else {
		// set the previous the base url the 1st time
		c.Previous = c.BaseApiUrl
		// fetch and encode from baseapiurl
		if err := api.FetchFromCache(c, c.BaseApiUrl, &ah); err != nil {
			c.Logger.Error("error doing fetching from cache", "error", err)
			return err
		}
		// set the next url
		c.Next = ah.Next

	}

	// loop through the results
	for _, l := range ah.Results {
		fmt.Println(l.Name)
	}

	return nil
}

// Mapb fetches the previous map of locations from the PokeAPI and displays them.
func Mapb(c *models.Config, args ...string) error {
	var ah models.Apiheader
	if c.Previous != "" {
		if err := api.FetchFromCache(c, c.Previous, &ah); err != nil {
			c.Logger.Error("error fetching from cache", "error", err)
			return err
		}
	} else {
		if err := api.FetchFromCache(c, c.BaseApiUrl, &ah); err != nil {
			c.Logger.Error("error fetching from cache", "error", err)
			return err
		}
	}

	// set the next url
	c.Next = ah.Next

	// loop through the results
	for _, l := range ah.Results {
		fmt.Println(l.Name)
	}

	return nil

}

// AltExplore is an alternative explore function that allows for more direct exploration
func Explore(c *models.Config, args ...string) error {
	// check if args are empty
	if args[0] == "" || args[1] == "" {
		c.Logger.Error("invalid location arguments", "args", args)
		return errors.New("invalid location")
	}

	// clean and trim them
	cleanLocation := strings.TrimSpace(strings.ToLower(args[1]))

	// check if previous url is set, if not set to base url
	if c.Previous == "" {
		c.Previous = c.BaseApiUrl
	}

	// download the json data and encode to struct
	var locationArea models.LocationArea
	locationUrl := fmt.Sprintf("%s%s", c.Previous, cleanLocation)
	api.FetchFromCache(c, locationUrl, &locationArea)

	fmt.Println("Found Pokemon:")
	for _, k := range locationArea.PokemonEncounters {
		fmt.Printf("- %s\n", k.Pokemon.Name)
	}
	return nil
}

// Catch attempts to catch a pokemon by throwing a Pokeball at it.
func Catch(c *models.Config, args ...string) error {
	if args[0] == "" || args[1] == "" {
		// log the error and return
		c.Logger.Error("invalid character arguments", "args", args)
		return errors.New("invalid character")
	}
	// fetch the character from args1
	character := args[1]
	fmt.Printf("Throwing a Pokeball at %s...\n", character)

	// fetch and encode the pokemon data from the api
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", character)
	var pokemon models.Pokemon
	api.FetchFromCache(c, url, &pokemon)

	// TODO: account for pokemon.BaseExp
	if randomChance(.50) {
		c.Pokedex[character] = pokemon

		// marshal json from caught pokemon
		jsonData, err := json.Marshal(pokemon)
		if err != nil {
			c.Logger.Error("failed to marshal pokemon", "error", err)
			return err
		}

		// write to database
		_, err = c.Db.AddPokemon(context.Background(), database.AddPokemonParams{
			ID:          uuid.New(),
			PokemonName: character,
			JsonData:    pqtype.NullRawMessage{RawMessage: jsonData, Valid: true},
		})
		if err != nil {
			c.Logger.Error("failed to insert pokemon", "error", err)
			return err
		}
		c.Logger.Info("pokemon caught", "pokemon", character)
		return nil
	}

	c.Logger.Info("pokemon escaped", "pokemon", character)
	return nil
}

// Inspect displays the details of a caught pokemon.
func Inspect(c *models.Config, args ...string) error {
	character := args[1]
	val, ok := c.Pokedex[character]
	if !ok {
		c.Logger.Error("pokemon not found", "pokemon", character)
		return nil
	}

	fmt.Printf("Name: %s\n", val.Name)
	fmt.Printf("Height: %d\n", val.Height)
	fmt.Printf("Weight: %d\n", val.Weight)
	fmt.Println("Stats:")
	for _, stat := range val.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range val.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
	return nil

}

// Pokedex displays the list of caught pokemon.
func Pokedex(c *models.Config, args ...string) error {
	fmt.Println("Your Pokedex:")
	if len(c.Pokedex) == 0 {
		fmt.Println("You have not caught any pokemon yet.")
		return nil
	}

	pokemon, err := c.Db.ListPokemon(context.Background())
	if err != nil {
		c.Logger.Error("error fetching pokemon from database", "error", err)
		return err
	}

	for _, b := range pokemon {
		fmt.Println(b.PokemonName)
	}
	return nil
}

// History shows all the commands that were previously used
func History(c *models.Config, args ...string) error {
	fmt.Println("History:")
	for _, c := range c.History {
		fmt.Printf("-%s\n", c)
	}
	return nil
}
