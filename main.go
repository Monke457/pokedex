package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	api "pokedex/internal/pkg/pokeapi"
)

var config = api.Config{}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()

	for {
		fmt.Printf("pokedex > ")

		scanner.Scan() 
		words := clean(scanner.Text())

		cmd, ok := commands[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if len(words) > 1 {
			config.Args = words[1:]
		}

		err := cmd.callback(&config)
		if err != nil {
			fmt.Println(err)
		}

		config.Args = []string{}
	}
}

func clean(val string) []string {
	val = strings.TrimSpace(val)
	val = strings.ToLower(val)
	words := strings.Fields(val)
	return words 
}
