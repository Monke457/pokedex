package main

import (
	"fmt"
	"os"

	api "pokedex/internal/pkg/pokeapi"
)

type cliCommand struct {
	name string
	description string
	callback func(*api.Config) error
}

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
			description: `Displays the names of 20 location areas 
			in the Pokemon world. Each subsequent call to map should 
			display the next 20 locations, and so on.`,
			callback: mapf,
		},
		"mapb": {
			name: "mapb",
			description: `Similar to the map command, however, 
			instead of displaying the next 20 locations, it displays 
			the previous 20 locations. It's a way to go back.`,
			callback: mapb,
		},
	}
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
