package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joshhartwig/pokedex/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type config struct {
	commands   map[string]cliCommand
	next       string
	previous   string
	baseApiUrl string
	cache      pokecache.Cache
}

// json decoding
type location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type apiheader struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []location `json:"results"`
}

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
				c.commands[cleanedInput[0]].callback()
			} else {
				fmt.Println("Uknown command")
			}
		}
	}
}

func (c *config) exitCmd() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func (c *config) helpCmd() error {
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

//TODO: Finish using the cache in the api instead of calling the api direct

func (c *config) mapCmd() error {

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
		// fetch and encode from baseapiurl
		if err := c.fetchFromCache(c.baseApiUrl, &ah); err != nil {
			fmt.Println("error doing fetch& encode on c.baseapiurl: ", c.baseApiUrl)
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

func (c *config) mapbCmd() error {
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

// func fetchAndEncode(url string, v any) error {
// 	// check cache 1st
// 	client := &http.Client{}

// 	resp, err := client.Get(url)
// 	if err != nil {
// 		return err
// 	}

// 	defer resp.Body.Close()

// 	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
// 		return err
// 	}
// 	return nil
// }

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

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		c.cache.Add(url, data)
		// return the data
		if err = json.NewDecoder(bytes.NewReader(data)).Decode(&v); err != nil {
			return err
		}
	}

	// we found the url in the cache, return the data
	if err := json.NewDecoder(bytes.NewReader(found.Val)).Decode(&v); err != nil {
		return err
	}

	return nil
}
