package main

import (
	"fmt"
	"math/rand"
	"os"

	api "pokedex/internal/pkg/pokeapi"
	dex "pokedex/internal/pkg/pokedex"
)

type cliCommand struct {
	name string
	description string
	callback func(*api.Config) error
}

var pokedex = dex.New()

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help" : {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"map": {
			name: "map",
			description: "Displays the names of 20 location areas " + 
			"in the Pokemon world. Each subsequent call to map should " +
			"display the next 20 locations, and so on.",
			callback: mapf,
		},
		"mapb": {
			name: "mapb",
			description: "Similar to the map command, however, " +
			"instead of displaying the next 20 locations, it displays " +
			"the previous 20 locations. It's a way to go back.",
			callback: mapb,
		},
		"explore": {
			name: "explore",
			description: "See a list of all the Pokémon in a given area. " +
			"Takes the name of a location as an argument",
			callback: explore,
		},
		"catch": {
			name: "catch",
			description: "Try to catch a pokemon. " +
			"Takes the name of a pokemon as an argument.",
			callback: catch,
		},
		"inspect": {
			name: "inspect",
			description: "Look up a Pokemon in your Pokedex. " +
			"Only works with Pokemon you have already caught.",
			callback: inspect,
		},
		"pokedex": {
			name: "pokedex",
			description: "List all the Pokemon you have recorded in your Pokedex",
			callback: list,
		},
	}
} 

func explore(config *api.Config) error {
	if len(config.Args) == 0 {
		return fmt.Errorf("You must provide a location to explore.")
	}

	fmt.Printf("Exploring %s...\n", config.Args[0])

	res, err := api.ExploreLocation(config.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Found Pokemon:\n")
	for _, pokemon := range(res) {
		fmt.Printf(" - %s\n", pokemon)
	}

	return nil
}

func catch(config *api.Config) error {
	if len(config.Args) == 0 {
		return fmt.Errorf("You must provide the name of a pokemon to catch.")
	}

	res, err := api.GetPokemon(config.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", config.Args[0])
	chance := rand.Intn(80 + res.BaseExperience)

	if chance < 100 {
		pokedex.Add(*res)
		fmt.Printf("%s was caught!\n", config.Args[0])
		fmt.Printf("You may now inspect it with the inspect command.\n")
	} else {
		fmt.Printf("%s escaped!\n", config.Args[0])
	}
	return nil
}

func inspect(config *api.Config) error {
	if len(config.Args) == 0 {
		return fmt.Errorf("You must provide the name of a pokemon to inspect.")
	}

	poke, ok := pokedex.Get(config.Args[0])
	if !ok {
		return fmt.Errorf("You do not have an entry for %s in your Pokedex.", config.Args[0]) 
	}

	poke.Print()

	return nil
}

func list(config *api.Config) error {
	for name := range pokedex.List() {
		fmt.Printf("  - %s\n", name)
	}
	return nil
}

func mapf(config *api.Config) error {
	res, err := api.GetLocations(config.Next)
	if err != nil {
		return err
	}
	for _, loc := range(res.Results) {
		fmt.Printf("%s\n", loc.Name)
	}
	config.Next = res.Next
	config.Previous = res.Previous
	return nil
}

func mapb(config *api.Config) error {
	res, err := api.GetLocations(config.Previous)
	if err != nil {
		return err
	}
	for _, loc := range(res.Results) {
		fmt.Printf("%s\n", loc.Name)
	}
	config.Next = res.Next
	config.Previous = res.Previous
	return nil
}

func commandHelp(config *api.Config) error {
	fmt.Printf("\nWelcome to the Pokedex!\n")
	fmt.Printf("Usage:\n\n")

	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()

	return nil
}

func commandExit(config *api.Config) error {
	os.Exit(1)
	return nil
}
