package repl

import (
	"math/rand"
	"strings"
)

// cleanInput takes a string and returns a slice of strings
// where each string is a word in the input, all in lowercase and trimmed of whitespace.
// This is useful for normalizing user input for commands.
func cleanInput(text string) []string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	return strings.Split(text, " ")
}

// randomChance takes a percentage as a float64 and returns true if a randomly generated float64
// is less than or equal to the percentage. This is useful for simulating random events
// such as catching a Pokémon or finding an item.
func randomChance(percentage float64) bool {
	random := rand.Float64()
	return random <= ((percentage * 100) / 100)
}
