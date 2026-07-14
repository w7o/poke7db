package main

/*
Custom unmarshalling logic for specific types found in types.go
*/

import (
	"encoding/json"
	"errors"
	"net/url"
	"path"
)

func (pokemon *Pokemon) UnmarshalJSON(data []byte) error {
	type Alias Pokemon

	var aux struct {
		Alias   // includes the rest of the struct
		Species struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	*pokemon = Pokemon(aux.Alias)

	// Grabbing last subdirectory of URL to get the dex number
	fullURL, err := url.Parse(aux.Species.URL)
	if err != nil {
		return err
	}

	pokemon.DexNum = path.Base(fullURL.Path)
	pokemon.Name = aux.Species.Name

	return nil
}

func (dexColor *DexColor) UnmarshalJSON(data []byte) error {
	var color struct {
		Name string `json:"name"`
	}

	err := json.Unmarshal(data, &color)
	if err != nil {
		return err
	}

	// Set DexColor to the color name stored in "color"
	*dexColor = DexColor(color.Name)
	return nil
}

func (statBlock *StatBlock) UnmarshalJSON(data []byte) error {
	var stats []struct {
		BaseStat int `json:"base_stat"`
		EVYield  int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	}

	// Unmarshal data to the API's stats struct
	err := json.Unmarshal(data, &stats)
	if err != nil {
		return err
	}

	// Maps the found string to a PokemonStat
	statMap := map[string]*PokemonStat{
		"hp":              &statBlock.HP,
		"attack":          &statBlock.Attack,
		"defense":         &statBlock.Defense,
		"special-attack":  &statBlock.SpecialAttack,
		"special-defense": &statBlock.SpecialDefense,
		"speed":           &statBlock.Speed,
	}

	for _, stat := range stats {
		target, ok := statMap[stat.Stat.Name]
		if !ok {
			return errors.New("Unmarshalling mapping error")
		}
		target.BaseStat = stat.BaseStat
		target.EVYield = stat.EVYield
	}
	return nil
}
