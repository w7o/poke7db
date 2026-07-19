package main

/*
Custom unmarshalling logic for specific types found in types.go
*/

import (
	"encoding/json"
	"errors"
	"net/url"
	"path"
	"strconv"
)

func (category *Category) UnmarshalJSON(data []byte) error {
	// Grabs specifically the English name of the Pokémon
	// Does not currently support other languages
	var genera []struct {
		Genus    string   `json:"genus"`
		Language Language `json:"language"`
	}

	err := json.Unmarshal(data, &genera)
	if err != nil {
		return err
	}

	for _, g := range genera {
		if g.Language.Name == "en" {
			*category = Category(g.Genus)
			return nil
		}
	}

	return errors.New("No English category entry found")
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

func (generation *Generation) UnmarshalJSON(data []byte) error {
	var g NamedResource

	err := json.Unmarshal(data, &g)
	if err != nil {
		return err
	}

	fullURL, err := url.Parse(g.URL)
	if err != nil {
		return err
	}

	// Converting the last subdirectory of the URL -- which should be an integer --
	// to a string
	gen, err := strconv.Atoi(path.Base(fullURL.Path))
	if err != nil {
		return err
	}

	*generation = Generation(gen)
	return nil
}

func (pokemon *Pokemon) UnmarshalJSON(data []byte) error {
	type Alias Pokemon

	var aux struct {
		Alias   // includes the rest of the struct
		Species struct {
			// Name string `json:"name"`
			URL string `json:"url"`
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
	// pokemon.Name = aux.Species.Name

	return nil
}

func (form *PokemonForm) UnmarshalJSON(data []byte) error {
	var variety struct {
		Pokemon NamedResource `json:"pokemon"`
	}

	err := json.Unmarshal(data, &variety)
	if err != nil {
		return err
	}

	fullURL, err := url.Parse(variety.Pokemon.URL)
	if err != nil {
		return err
	}

	n, err := strconv.Atoi(path.Base(fullURL.Path))
	if err != nil {
		return err
	}

	form.Name = variety.Pokemon.Name
	form.APINumber = n

	return nil
}

func (name *PokemonName) UnmarshalJSON(data []byte) error {
	var langNames []struct {
		Language Language `json:"language"`
		Name     string   `json:"name"`
	}

	err := json.Unmarshal(data, &langNames)
	if err != nil {
		return err
	}

	for _, n := range langNames {
		if n.Language.Name == "en" {
			*name = PokemonName(n.Name)
			return nil
		}
	}
	return errors.New("No English PokemonName entry found")
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
			return errors.New("Stat block unmarshalling mapping error")
		}
		target.BaseStat = stat.BaseStat
		target.EVYield = stat.EVYield
	}
	return nil
}
