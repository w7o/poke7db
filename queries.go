package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func json_query(link string) ([]byte, error) {
	log.SetPrefix("json_query: ")
	log.SetFlags(0)

	//Empty string check
	if link == "" {
		return nil, errors.New("Empty endpoint")
	}

	//Send request to endpoint to get a response
	log.Printf("Request to %v", link)
	response, err := http.Get(link)

	//Error check
	if err != nil {
		return nil, fmt.Errorf("request %q failed: %w", link, err)
	}

	//Close request
	defer response.Body.Close()

	//Read the response's data
	data, err := io.ReadAll(response.Body)

	//Error check
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	//Add to localdata if dev is not enabled
	if os.Getenv("P7D_ENV") != "dev" {
		u, err := url.Parse(link)
		if err != nil {
			return nil, fmt.Errorf("storage to localdata failed: %w", err)
		}

		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		parts = parts[2:]
		parts[len(parts)-1] += ".json"

		// "..." passes each element in the list as a separate argument
		filename := filepath.Join("localdata", filepath.Join(parts...))

		// file permission written in octal (rwxrwxrwx)
		err = os.MkdirAll(filepath.Dir(filename), 0o777)
		if err != nil {
			return nil, fmt.Errorf("creating directory in localdata failed: %w", err)
		}

		err = os.WriteFile(filename, data, 0o666)
		if err != nil {
			return nil, fmt.Errorf("writing to localdata failed: %w", err)
		}
	}

	return data, nil
}

/*
 */
func merge_pokemon_structs(api_species *APIPokemonSpecies, pokemon *Pokemon) {
	pokemon.SpeciesInfo = api_species
}

/*
Queries the API for Pokémon statistics, and
returns a struct containing all exported information
*/
func poke_query(id string) (Pokemon, error) {
	// POKEMON
	url := SITE + "pokemon/" + id
	// Grab JSON data using url
	data, err := json_query(url)

	if err != nil {
		return Pokemon{}, fmt.Errorf("Pokemon/ JSON query failed: %w", err)
	}

	var pokemon Pokemon
	err = json.Unmarshal(data, &pokemon)

	if err != nil {
		return Pokemon{}, fmt.Errorf("JSON unmarshal for Pokémon %q failed: %w", id, err)
	}

	// POKEMON-SPECIES
	url = SITE + "pokemon-species/" + strconv.Itoa(pokemon.DexNum)
	data, err = json_query(url)

	if err != nil {
		return Pokemon{}, fmt.Errorf("Pokemon-species/ JSON query failed: %w", err)
	}

	var api_species APIPokemonSpecies
	err = json.Unmarshal(data, &api_species)
	// fmt.Println(api_species.FormsList)

	if err != nil {
		return Pokemon{}, fmt.Errorf("JSON unmarshal for Pokémon Species %q failed: %w", id, err)
	}

	merge_pokemon_structs(&api_species, &pokemon)

	return pokemon, nil
}
