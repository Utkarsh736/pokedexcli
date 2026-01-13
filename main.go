package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Utkarsh736/pokedexcli/internal/pokecache"
	"github.com/Utkarsh736/pokedexcli/internal/pokeapi"
)

type config struct {
	nextLocationURL     *string
	previousLocationURL *string
	cache               pokecache.Cache 
	caughtPokemon       map[string]pokeapi.Pokemon
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())


	scanner := bufio.NewScanner(os.Stdin)
	
	// Initialize config
	cfg := &config{
		cache: pokecache.NewCache(5 * time.Minute),
		caughtPokemon: make(map[string]pokeapi.Pokemon),
	}
	
	commands := getCommands()
	
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		
		words := cleanInput(input)
		
		if len(words) == 0 {
			continue
		}
		
		commandName := words[0]
		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}


		command, exists := commands[commandName]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}
		
		// Pass config to the command callback - THIS WAS THE MISSING PART
		err := command.callback(cfg, args...)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

