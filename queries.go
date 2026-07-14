package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	// "encoding/json"
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
Queries the API and returns a struct containing all exported information
*/
func poke_query(id string) (PokemonSpecies, error) {
	endpoint := SITE + "pokemon/" + id
	// Grab JSON data using endpoint
	data, err := json_query(endpoint)

	if err != nil {
		return PokemonSpecies{}, fmt.Errorf("JSON query failed: %w", err)
	}

	// Unpack the JSON data into a PokemonSpecies struct
	var pokemon PokemonSpecies
	json.Unmarshal(data, &pokemon)

	if err != nil {
		return PokemonSpecies{}, fmt.Errorf("JSON unmarshal for Pokémon %q failed: %w", id, err)
	}

	return pokemon, nil
}
