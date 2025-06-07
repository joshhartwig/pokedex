package repl

import (
	"bufio"
	"fmt"
	"os"

	"github.com/joshhartwig/pokedex/pkg/models"
)

// Repl starts a Read-Eval-Print Loop for the Pokedex application.
// It reads user input from the command line, processes commands, and executes the corresponding callbacks.
// It will continue to prompt for input until the user exits the application.
func Repl(c *models.Config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("Pokedex > ")
		if scanner.Scan() {
			input := scanner.Text()
			cleanedInput := cleanInput(input)

			// user must enter a command
			if len(cleanedInput) == 0 {
				fmt.Println("Please enter a command. Type help for assistance.")
				continue
			}

			cmd, ok := c.Commands[cleanedInput[0]]
			if !ok {
				fmt.Println("Uknown Command, type Help for assistance")
				continue
			}

			// add the command to the history
			c.History = append(c.History, cleanedInput[0])

			err := cmd.Callback(cleanedInput...)
			if err != nil {
				fmt.Println("The command you just ran requires additional requirements")
			}
		}
	}
}
