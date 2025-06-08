package repl

import (
	"errors"
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

// CatchPokemon attempts to catch a Pokemon based on a base catch chance and the Pokemon's base experience.
// The final catch probability is calculated by reducing the base chance according to the Pokemon's experience level,
// with higher-level Pokemon being harder to catch. The final probability is clamped between 1% and 95%.
//
// Parameters:
//   - baseChance: The initial probability of catching the Pokemon (between 0.0 and 1.0)
//   - baseExperience: The Pokemon's base experience value, affecting catch difficulty
//
// Returns:
//   - bool: true if the Pokemon was caught, false otherwise
func CatchPokemon(baseChance float64, baseExperience int) bool {
	const experienceFactor = 0.001 // tuning constant — adjust to your liking

	// Compute final catch chance
	catchChance := baseChance - (float64(baseExperience) * experienceFactor)

	// Clamp between 1% and 95% to avoid extremes
	if catchChance < 0.01 {
		catchChance = 0.01
	}
	if catchChance > 0.95 {
		catchChance = 0.95
	}

	// Perform the catch roll
	random := rand.Float64()
	return random <= catchChance
}

// checkArgs validates command line arguments.
// It ensures there are at least two arguments provided and neither is empty.
// Returns an error if validation fails, nil otherwise.
// args: Slice of strings containing command line arguments
func checkArgs(wantedLength int, args []string) error {
	if len(args) < 2 {
		return errors.New("this command requires two arguments")
	}

	for x := 0; x < wantedLength; x++ {
		if strings.TrimSpace(args[x]) == "" {
			return errors.New("arguments cannot be empty")
		}
	}

	return nil
}
