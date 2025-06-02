package repl

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joshhartwig/pokedex/internal/api"
	"github.com/joshhartwig/pokedex/pkg/models"
)

func Exit(c *models.Config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func Help(c *models.Config, args ...string) error {
	msg := `
	Welcome to the Pokedex!\n
	Usage:\n\n
	help: Displays a help message\n
	map: Displays map of locations from api\n
	mapb: Displays the previous locations from the api\n
	exit: Exit the pokedex\n
	`
	fmt.Printf("%s", msg)
	return nil
}

func Map(c *models.Config, args ...string) error {

	var ah models.Apiheader
	// if c.next is anything but empty, it likely has a url and pull from that url
	if c.Next != "" {
		fmt.Println("c.next has value: ", c.Next)
		if err := api.FetchFromCache(c, c.Next, &ah); err != nil {
			fmt.Println("error fetching c.next")
			return err
		}
	} else {
		fmt.Println("c.next did not have value, go to baserul")
		// set the previous the base url the 1st time
		c.Previous = c.BaseApiUrl
		// fetch and encode from baseapiurl
		if err := api.FetchFromCache(c, c.BaseApiUrl, &ah); err != nil {
			fmt.Println("error doing fetch& encode on c.baseapiurl: ", c.BaseApiUrl, err)
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

func Mapb(c *models.Config, args ...string) error {
	var ah models.Apiheader
	if c.Previous != "" {
		if err := api.FetchFromCache(c, c.Previous, &ah); err != nil {
			return err
		}
	} else {
		if err := api.FetchFromCache(c, c.BaseApiUrl, &ah); err != nil {
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

func Explore(c *models.Config, args ...string) error {
	if args[0] == "" || args[1] == "" {
		return errors.New("invalid location")
	}

	cleanLocation := strings.TrimSpace(strings.ToLower(args[1]))

	// fetch the location api data to get the location names, encode to apiheader struct and loop
	// through them to determine if there is a match with our explore - name
	var ah models.Apiheader
	err := api.FetchFromCache(c, c.Previous, &ah)
	if err != nil {
		fmt.Println("error fetching from cache:", err)
		return err
	}

	// loop through apiheader.results and find the location name
	for _, v := range ah.Results {
		if v.Name == cleanLocation {
			locationUrl := fmt.Sprintf("%s%s", c.Previous, cleanLocation)
			var locationArea models.LocationArea
			// fetch from cache or download and encode to location area struct
			api.FetchFromCache(c, locationUrl, &locationArea)

			// loop through the pokemon encounters and list the pokemon names
			fmt.Println("Found Pokemon:")
			for _, k := range locationArea.PokemonEncounters {
				fmt.Printf("- %s\n", k.Pokemon.Name)
			}
		}
	}
	return nil
}

func AltExplore(c *models.Config, args ...string) error {
	// check if args are empty
	if args[0] == "" || args[1] == "" {
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

func Catch(c *models.Config, args ...string) error {
	if args[0] == "" || args[1] == "" {
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
		c.Pokedex[character] += pokemon
		fmt.Printf("%s was caught!\n", character)
		return nil
	}

	fmt.Printf("%s escaped!\n", character)
	return nil
}

func Inspect(c *models.Config, args ...string) error {
	character := args[1]
	val, ok := c.Pokedex[character]
	if !ok {
		fmt.Println("you have not caught that pokemon")
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

/*
Add an inspect command. It takes the name of a Pokemon and prints the name, height, weight, stats and type(s) of the Pokemon. Example usage:

Pokedex > inspect pidgey
you have not caught that pokemon
Pokedex > catch pidgey
Throwing a Pokeball at pidgey...
pidgey was caught!
Pokedex > inspect pidgey
Name: pidgey
Height: 3
Weight: 18
Stats:
  -hp: 40
  -attack: 45
  -defense: 40
  -special-attack: 35
  -special-defense: 35
  -speed: 56
Types:
  - normal
  - flying
*/
