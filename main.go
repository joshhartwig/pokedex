package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
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

// TODO: map is not working now that i changed the flow, but mapb was workign
func (c *config) mapCmd() error {

	var ah apiheader
	if c.next != "" {
		err := fetchAndEncode(c.next, &ah)
		if err != nil {
			fmt.Println("error fetching c.next")
			return err
		} else {
			err := fetchAndEncode(c.baseApiUrl, &ah)
			if err != nil {
				fmt.Println("error fetching c.baseapiurl", c.baseApiUrl)
				return err
			}
		}

		// set the next url
		c.next = ah.Next

		// loop through the results
		for _, l := range ah.Results {
			fmt.Println(l.Name)
		}

	}
	return nil
}

func (c *config) mapbCmd() error {
	var ah apiheader
	if c.previous != "" {
		if err := fetchAndEncode(c.previous, &ah); err != nil {
			return err
		}
	} else {
		if err := fetchAndEncode(c.baseApiUrl, &ah); err != nil {
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

func fetchAndEncode(url string, v any) error {
	client := &http.Client{}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return err
	}
	return nil
}
