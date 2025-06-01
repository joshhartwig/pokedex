package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joshhartwig/pokedex/internal/pokecache"
)

func main() {
	app := config{}
	app.commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "used to close the app",
			callback:    app.exitCmd,
		},
		"help": {
			name:        "help",
			description: "used to get help",
			callback:    app.helpCmd,
		},
		"map": {
			name:        "map",
			description: "used to list all the pokedex locations",
			callback:    app.mapCmd,
		},
		"mapb": {
			name:        "mapb",
			description: "used to move forward in the map",
			callback:    app.mapbCmd,
		},
		"explore": {
			name:        "explore",
			description: "explores a section of the map",
			callback:    app.altExploreCmd,
		},
	}
	app.baseApiUrl = "https://pokeapi.co/api/v2/location-area/"

	// create a new cache with a timer of 10 seconds
	app.cache = *pokecache.NewCache(time.Millisecond * 10)

	// start the repl loop
	app.repl()
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	return strings.Split(text, " ")
}

func (c *config) repl() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("Pokedex > ")
		if scanner.Scan() {
			input := scanner.Text()
			cleanedInput := cleanInput(input)
			_, ok := c.commands[cleanedInput[0]]
			if ok {
				//TODO: bug where if there is no cleanedInput[1] it will crash
				err := c.commands[cleanedInput[0]].callback(cleanedInput...)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Uknown command")
			}
		}
	}
}

func (c *config) exitCmd(s ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func (c *config) helpCmd(s ...string) error {
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

func (c *config) mapCmd(s ...string) error {

	var ah apiheader
	// if c.next is anything but empty, it likely has a url and pull from that url
	if c.next != "" {
		fmt.Println("c.next has value: ", c.next)
		if err := c.fetchFromCache(c.next, &ah); err != nil {
			fmt.Println("error fetching c.next")
			return err
		}
	} else {
		fmt.Println("c.next did not have value, go to baserul")
		// set the previous the base url the 1st time
		c.previous = c.baseApiUrl
		// fetch and encode from baseapiurl
		if err := c.fetchFromCache(c.baseApiUrl, &ah); err != nil {
			fmt.Println("error doing fetch& encode on c.baseapiurl: ", c.baseApiUrl, err)
			return err
		}
		// set the next url
		c.next = ah.Next

	}

	// loop through the results
	for _, l := range ah.Results {
		fmt.Println(l.Name)
	}

	return nil
}

func (c *config) mapbCmd(s ...string) error {
	var ah apiheader
	if c.previous != "" {
		if err := c.fetchFromCache(c.previous, &ah); err != nil {
			return err
		}
	} else {
		if err := c.fetchFromCache(c.baseApiUrl, &ah); err != nil {
			return err
		}
	}

	// set the next url
	c.next = ah.Next

	// loop through the results
	for _, l := range ah.Results {
		fmt.Println(l.Name)
	}

	return nil

}

// fetchFromCache looks for url string in cache and returns associated data
// if the url is not found in the cache it will download the data and add the
// url and associated data encoded in json format to the map and return it
func (c *config) fetchFromCache(url string, v any) error {
	// try to find the url in cache 1st
	found, ok := c.cache.Entries[url]
	if !ok { // if we did not find it, add a new cache entry with the data
		client := &http.Client{}
		resp, err := client.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// convert the resp.body to byte slide
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// add to cache
		c.cache.Add(url, data)

		// immediately fetch the data from the cache
		b, _ := c.cache.Get(url)
		if err := json.NewDecoder(bytes.NewReader(b)).Decode(&v); err != nil {
			fmt.Println("error decoding bytes", err)
			return err
		}
		return nil
	}

	// we found the url in the cache, return the data
	if err := json.NewDecoder(bytes.NewReader(found.Val)).Decode(&v); err != nil {
		fmt.Println("error decoding:", err)
		return err
	}

	return nil
}

func (c *config) exploreCmd(args ...string) error {
	if args[0] == "" || args[1] == "" {
		return errors.New("invalid location")
	}

	cleanLocation := strings.TrimSpace(strings.ToLower(args[1]))

	// fetch the location api data to get the location names, encode to apiheader struct and loop
	// through them to determine if there is a match with our explore - name
	var ah apiheader
	err := c.fetchFromCache(c.previous, &ah)
	if err != nil {
		fmt.Println("error fetching from cache:", err)
		return err
	}

	// loop through apiheader.results and find the location name
	for _, v := range ah.Results {
		if v.Name == cleanLocation {
			locationUrl := fmt.Sprintf("%s%s", c.previous, cleanLocation)
			var locationArea LocationArea
			// fetch from cache or download and encode to location area struct
			c.fetchFromCache(locationUrl, &locationArea)

			// loop through the pokemon encounters and list the pokemon names
			fmt.Println("Found Pokemon:")
			for _, k := range locationArea.PokemonEncounters {
				fmt.Printf("- %s\n", k.Pokemon.Name)
			}
		}
	}
	return nil
}

func (c *config) altExploreCmd(args ...string) error {
	// check if args are empty
	if args[0] == "" || args[1] == "" {
		return errors.New("invalid location")
	}

	// clean and trim them
	cleanLocation := strings.TrimSpace(strings.ToLower(args[1]))

	// check if previous url is set, if not set to base url
	if c.previous == "" {
		c.previous = c.baseApiUrl
	}

	// download the json data and encode to struct
	var locationArea LocationArea
	locationUrl := fmt.Sprintf("%s%s", c.previous, cleanLocation)
	c.fetchFromCache(locationUrl, &locationArea)

	fmt.Println("Found Pokemon:")
	for _, k := range locationArea.PokemonEncounters {
		fmt.Printf("- %s\n", k.Pokemon.Name)
	}
	return nil
}
