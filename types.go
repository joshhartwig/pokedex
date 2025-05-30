package main

import (
	"github.com/joshhartwig/pokedex/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(...string) error
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
