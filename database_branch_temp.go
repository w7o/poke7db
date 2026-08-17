package main

import (
	"database/sql"
	"fmt"
	"reflect"
	//"strconv"
	"strings"
)

/*
Temporary branch because I accidentally left my v0.1.4 - v0.1.5 commits unpushed
*/

/*
Struct type names MUST be in the format {TableName}TE
*/

// T01/00
type PokemonSpeciesTE struct {
	PokemonSpeciesID string // alphanum dex numebrs
	PokemonName      string
	Category         string
	BaseHappiness    int
	CaptureRate      int
	GrowthRateID     int
	GenderRate       int
	HatchCounter     int
	ColorID          int
	ShapeID          int
	IsMythical       int
	IsLegendary      int
}

var pokemonSpecies []PokemonSpeciesTE

// T01/01
type PokemonFormTE struct {
	PokemonFormID       int
	PokemonSpeciesID    string
	FormName            string
	StatHP              int
	StatDefense         int
	StatAttack          int
	StatSpecialAttack   int
	StatSpecialDefense  int
	StatSpeed           int
	Height              int
	Weight              int
	BaseExperienceYield int
}

var pokemonForm []PokemonFormTE

// T01/02
type PokemonTypeTE struct {
	PokemonFormID int
	Slot          int
	TypeID        int
}

var pokemonTypes []PokemonTypeTE

// T01/03
type PokemonEggGroupTE struct {
	PokemonSpeciesID string
	Slot             int
	EggGroupID       int
}

var pokemonEggGroup []PokemonEggGroupTE

// // T01/04
// type PokemonEVYieldTE struct {
// 	PokemonSpeciesID string
// 	StatID           int
// 	EVYield          int
// }

// var pokemonEVYields []PokemonEVYieldTE

func insertStruct[T any](tx *sql.Tx, tableName string, dbStruct []T) error {
	// Grab the type of the struct
	structType := reflect.TypeOf(dbStruct[0])

	// Grab the value of the structs
	structValue := reflect.ValueOf(dbStruct)

	// Check if T is struct
	if structType.Kind() != reflect.Struct {
		return retError("F_03", "Passed in type is not a struct", nil)
	}

	// Returns the number of fields in the struct
	fields := structType.NumField()

	// // grabs table name from {TableName}TE
	// tableName := strings.TrimSuffix(structType.Name(), "TE")
	// if tableName != structType.Name() {
	// 	return retError("F_02", "Struct name doesn't have TE suffix", nil)
	// }

	// if need DB tag implementation: sDC 719-06#4

	var fieldNames []string
	var fieldValues []any
	var placeholders []string
	for i := 0; i < fields; i++ {
		fieldNames = append(fieldNames, structType.Field(i).Name)
	}

	var placeholderIndex int = 1
	for i := 0; i < structValue.Len(); i++ {
		currStruct := structValue.Index(i)
		var rowPlaceholders []string
		for j := 0; j < fields; j++ {
			fieldValues = append(fieldValues, currStruct.Field(j).Interface())
			rowPlaceholders = append(rowPlaceholders,
				fmt.Sprintf("$%d", placeholderIndex))
			placeholderIndex++
		}
		placeholders = append(placeholders,
			"("+strings.Join(rowPlaceholders, ", ")+")",
		)
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		tableName,
		strings.Join(fieldNames, ", "),
		strings.Join(placeholders, ", "),
	)

	_, err := tx.Exec(query, fieldValues...)
	if err != nil {
		writeLog(query)
		writeLog(fmt.Sprintf("values: %v", fieldValues))
		return retError("F_04", "Execution of query failed; offending query saved in log.txt", err)
	}

	return nil
}

func initTemporaryData() error {
	pokemonSpecies = append(pokemonSpecies, PokemonSpeciesTE{
		PokemonSpeciesID: "197",
		PokemonName:      "Umbreon",
		Category:         "Moonlight Pokémon",
		BaseHappiness:    35,
		CaptureRate:      45,
		GrowthRateID:     0,
		GenderRate:       2,
		HatchCounter:     35,
		ColorID:          4,
		ShapeID:          0,
		IsMythical:       0,
		IsLegendary:      0,
	})
	pokemonForm = append(pokemonForm, PokemonFormTE{
		PokemonFormID:       197,
		PokemonSpeciesID:    "197",
		FormName:            "Umbreon",
		StatHP:              95,
		StatAttack:          65,
		StatDefense:         110,
		StatSpecialAttack:   60,
		StatSpecialDefense:  130,
		StatSpeed:           65,
		Height:              10,
		Weight:              270,
		BaseExperienceYield: 184,
	})
	pokemonTypes = append(pokemonTypes,
		PokemonTypeTE{
			PokemonFormID: 197,
			Slot:          0, //%TODO remember to make this 0-indexed when importing from pokeapi
			TypeID:        17,
		})
	pokemonEggGroup = append(pokemonEggGroup,
		PokemonEggGroupTE{
			PokemonSpeciesID: "197",
			Slot:             0,
			EggGroupID:       8,
		})
	// pokemonEVYields = append(pokemonEVYields,
	// 	PokemonEVYieldTE{
	// 		PokemonSpeciesID: "197",
	// 		StatID:           4,
	// 		EVYield:          2,
	// 	})

	tx, err := database.Begin()
	if err != nil {
		return retError("F_00", "", err)
	}
	defer tx.Rollback()

	type table struct {
		TableName string
		tableVar  any
	}

	var tables []table
	tables = append(tables,
		table{"PokemonSpecies", pokemonSpecies},
		table{"PokemonForms", pokemonForm},
	)

	err = insertStruct(tx, "PokemonSpecies", pokemonSpecies)
	if err != nil {
		return err
	}

	err = insertStruct(tx, "PokemonForms", pokemonForm)
	if err != nil {
		return err
	}

	return nil
}
