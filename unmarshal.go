package main

/*
Custom unmarshalling logic for specific types found in types.go
*/

import (
	"encoding/json"
	"errors"
)

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
