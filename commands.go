package main

import (
	"fmt"
	"os"
	"errors"
	"math/rand"

	"github.com/Utkarsh736/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display previous 20 location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "View details about a caught Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "View all caught Pokemon",
			callback:    commandPokedex,
		},
	}
}


func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	
	commands := getCommands()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	
	return nil
}


func commandMap(cfg *config, args ...string) error {
	// Get location areas using the next URL
	resp, err := pokeapi.GetLocationAreas(cfg.nextLocationURL, &cfg.cache)
	if err != nil {
		return err
	}
	
	// Update config with new URLs
	cfg.nextLocationURL = resp.Next
	cfg.previousLocationURL = resp.Previous
	
	// Print location names
	for _, location := range resp.Results {
		fmt.Println(location.Name)
	}
	
	return nil
}

func commandMapb(cfg *config, args ...string) error {
	// Check if we're on the first page
	if cfg.previousLocationURL == nil {
		return errors.New("you're on the first page")
	}
	
	// Get location areas using the previous URL
	resp, err := pokeapi.GetLocationAreas(cfg.previousLocationURL, &cfg.cache)
	if err != nil {
		return err
	}
	
	// Update config with new URLs
	cfg.nextLocationURL = resp.Next
	cfg.previousLocationURL = resp.Previous
	
	// Print location names
	for _, location := range resp.Results {
		fmt.Println(location.Name)
	}
	
	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("you must provide a location area name")
	}

	locationAreaName := args[0]

	fmt.Printf("Exploring %s...\n", locationAreaName)

	// Get detailed location area data
	locationArea, err := pokeapi.GetLocationArea(locationAreaName, &cfg.cache)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range locationArea.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("you must provide a pokemon name")
	}

	pokemonName := args[0]

	// Check if already caught
	if _, exists := cfg.caughtPokemon[pokemonName]; exists {
		fmt.Printf("You have already caught %s!\n", pokemonName)
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	// Get Pokemon data
	pokemon, err := pokeapi.GetPokemon(pokemonName, &cfg.cache)
	if err != nil {
		return err
	}

	// Calculate catch chance based on base experience
	// Higher base experience = harder to catch
	const maxBaseExperience = 300
	catchThreshold := maxBaseExperience - pokemon.BaseExperience

	// Generate random number between 0 and maxBaseExperience
	randomValue := rand.Intn(maxBaseExperience)

	if randomValue < catchThreshold {
		// Caught!
		cfg.caughtPokemon[pokemonName] = pokemon
		fmt.Printf("%s was caught!\n", pokemonName)
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		// Escaped
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("you must provide a pokemon name")
	}

	pokemonName := args[0]

	// Check if the Pokemon has been caught
	pokemon, exists := cfg.caughtPokemon[pokemonName]
	if !exists {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	// Display Pokemon information
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Printf("  - %s\n", typeInfo.Type.Name)
	}

	return nil
}

func commandPokedex(cfg *config, args ...string) error {
	if len(cfg.caughtPokemon) == 0 {
		fmt.Println("You haven't caught any pokemon yet!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range cfg.caughtPokemon {
		fmt.Printf("  - %s\n", name)
	}

	return nil
}


