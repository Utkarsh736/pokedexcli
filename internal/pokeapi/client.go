package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	
	"github.com/Utkarsh736/pokedexcli/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2"

type LocationAreasResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationArea struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	PokemonEncounters    []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pokemon"`
}

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func GetPokemon(pokemonName string, cache *pokecache.Cache) (Pokemon, error) {
	url := baseURL + "/pokemon/" + pokemonName

	// Check cache first
	if cachedData, ok := cache.Get(url); ok {
		var pokemon Pokemon
		err := json.Unmarshal(cachedData, &pokemon)
		if err != nil {
			return Pokemon{}, err
		}
		return pokemon, nil
	}

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}

	// Add to cache
	cache.Add(url, body)

	var pokemon Pokemon
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return Pokemon{}, err
	}

	return pokemon, nil
}

func GetLocationArea(locationAreaName string, cache *pokecache.Cache) (LocationArea, error) {
	url := baseURL + "/location-area/" + locationAreaName

	// Check cache first
	if cachedData, ok := cache.Get(url); ok {
		var locationArea LocationArea
		err := json.Unmarshal(cachedData, &locationArea)
		if err != nil {
			return LocationArea{}, err
		}
		return locationArea, nil
	}

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return LocationArea{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationArea{}, err
	}

	// Add to cache
	cache.Add(url, body)

	var locationArea LocationArea
	err = json.Unmarshal(body, &locationArea)
	if err != nil {
		return LocationArea{}, err
	}

	return locationArea, nil
}


func GetLocationAreas(pageURL *string, cache *pokecache.Cache) (LocationAreasResponse, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}
	
	// Check if we have cached data
	if cachedData, ok := cache.Get(url); ok {
		// Cache hit! Unmarshal and return
		var locationResp LocationAreasResponse
		err := json.Unmarshal(cachedData, &locationResp)
		if err != nil {
			return LocationAreasResponse{}, err
		}
		return locationResp, nil
	}
	
	// Cache miss - make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return LocationAreasResponse{}, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreasResponse{}, err
	}
	
	// Add to cache for next time
	cache.Add(url, body)
	
	var locationResp LocationAreasResponse
	err = json.Unmarshal(body, &locationResp)
	if err != nil {
		return LocationAreasResponse{}, err
	}
	
	return locationResp, nil
}

