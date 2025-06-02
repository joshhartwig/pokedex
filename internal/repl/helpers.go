package repl

import (
	"math/rand"
	"strings"
)

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	return strings.Split(text, " ")
}

// returns true percentage of the time
func randomChance(percentage float64) bool {
	random := rand.Float64()
	return random <= ((percentage * 100) / 100)
}
