package repl

import (
	"bufio"
	"fmt"
	"os"

	"github.com/joshhartwig/pokedex/pkg/models"
)

func Repl(c *models.Config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("Pokedex > ")
		if scanner.Scan() {
			input := scanner.Text()
			cleanedInput := cleanInput(input)
			_, ok := c.Commands[cleanedInput[0]]
			if ok {
				//TODO: bug where if there is no cleanedInput[1] it will crash
				err := c.Commands[cleanedInput[0]].Callback(cleanedInput...)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Uknown command")
			}
		}
	}
}
