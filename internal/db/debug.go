package db

import (
	"database/sql"

	"github.com/w7o/poke7db/internal/logging"
)

var d_dexNum string

var d_testQueries = map[string]string{
	"PokemonSpecies": "SELECT * FROM PokemonSpecies WHERE PokemonSpeciesID = ",

	"PokemonForm": "SELECT * FROM PokemonForm WHERE PokemonFormID = ",

	"PokemonType": "SELECT * FROM PokemonType WHERE PokemonFormID = ",

	"PokemonEggGroup": "SELECT * FROM PokemonEggGroup WHERE PokemonSpeciesID = ",

	"PokemonEVYield": "SELECT * FROM PokemonEVYield WHERE PokemonFormID = ",
}

func d_queryAndLog(tx *sql.Tx, query string) error {
	r, err := tx.Query(query)
	if err != nil {
		return err
	}
	rr, err := StringRows(r)
	if err != nil {
		return err
	}
	logging.WriteLog(rr)
	return nil
}

func D_setDexNum(num string) {
	d_dexNum = num
}
