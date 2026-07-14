package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func json_query(endpoint string) ([]byte, error) {
	log.SetPrefix("json_query: ")
	log.SetFlags(0)

	//Empty string check
	if endpoint == "" {
		return nil, errors.New("Empty endpoint")
	}

	//Send request to endpoint to get a response
	log.Printf("Request to %v", endpoint)
	response, err := http.Get(endpoint)

	//Error check
	if err != nil {
		return nil, fmt.Errorf("request %q failed: %w", endpoint, err)
	}

	//Close request
	defer response.Body.Close()

	//Read the response's data
	data, err := io.ReadAll(response.Body)

	//Error check
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	return data, nil
}

/*
 */
func merge_pokemon_structs(api_species *APIPokemonSpecies, pokemon *Pokemon) {
	pokemon.BHappiness = api_species.BHappiness
	pokemon.CaptureRate = api_species.CaptureRate
	pokemon.Color = api_species.Color
	pokemon.EggGroups = api_species.EggGroups
}

/*
Queries the API for Pokémon statistics, and
returns a struct containing all exported information
*/
func poke_query(id string) (Pokemon, error) {
	// POKEMON
	endpoint := SITE + "pokemon/" + id
	// Grab JSON data using endpoint
	data, err := json_query(endpoint)

	if err != nil {
		return Pokemon{}, fmt.Errorf("Pokemon/ JSON query failed: %w", err)
	}

	var pokemon Pokemon
	err = json.Unmarshal(data, &pokemon)

	if err != nil {
		return Pokemon{}, fmt.Errorf("JSON unmarshal for Pokémon %q failed: %w", id, err)
	}

	// POKEMON-SPECIES
	endpoint = SITE + "pokemon-species/" + id
	data, err = json_query(endpoint)

	if err != nil {
		return Pokemon{}, fmt.Errorf("Pokemon-species/ JSON query failed: %w", err)
	}

	var api_species APIPokemonSpecies
	err = json.Unmarshal(data, &api_species)

	if err != nil {
		return Pokemon{}, fmt.Errorf("JSON unmarshal for Pokémon Species %q failed: %w", id, err)
	}

	merge_pokemon_structs(&api_species, &pokemon)

	return pokemon, nil
}
