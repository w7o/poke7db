package main

import (
	"database/sql"
	"fmt"
	"reflect"
	//"strconv"
	"strings"
)

var pokemonSpecies []PokemonSpeciesDBEntry
var pokemonForm []PokemonFormDBEntry
var pokemonType []PokemonTypeDBEntry
var pokemonEggGroup []PokemonEggGroupDBEntry
var pokemonEVYields []PokemonEVYieldDBEntry

func insertStruct(tx *sql.Tx, tableName string, dbStruct any) error {

	// Grab the value of the structs
	structValue := reflect.ValueOf(dbStruct)

	if structValue.Kind() != reflect.Slice {
		return retError("F_01", "Passed in type is not a list", nil)
	}

	if structValue.Len() == 0 {
		return retError("F_02", "Passed in list doesn't contain anything", nil)
	}

	// structvalue[0] doesnt work
	structType := structValue.Index(0).Type()

	// Check if T is struct
	if structValue.Len() == 0 {
		return retError("F_03", "Passed in type is not a struct", nil)
	}

	// Returns the number of fields in the struct
	fields := structType.NumField()

	// if need DB tag implementation: sDC 719-06#4

	var fieldNames []string
	var fieldValues []any
	var placeholders []string
	for i := range fields {
		fieldNames = append(fieldNames, structType.Field(i).Name)
	}

	var placeholderIndex int = 1
	for i := 0; i < structValue.Len(); i++ {
		currStruct := structValue.Index(i)
		var rowPlaceholders []string
		for j := range fields {
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
	pokemonSpecies = append(pokemonSpecies, PokemonSpeciesDBEntry{
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
	pokemonForm = append(pokemonForm, PokemonFormDBEntry{
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
	pokemonType = append(pokemonType,
		PokemonTypeDBEntry{
			PokemonFormID: 197,
			Slot:          0, //%TODO remember to make this 0-indexed when importing from pokeapi
			TypeID:        17,
		})
	pokemonEggGroup = append(pokemonEggGroup,
		PokemonEggGroupDBEntry{
			PokemonSpeciesID: "197",
			Slot:             0,
			EggGroupID:       8,
		})
	pokemonEVYields = append(pokemonEVYields,
		PokemonEVYieldDBEntry{
			PokemonFormID: 197,
			StatID:        4,
			EVYield:       2,
		})

	tx, err := database.Begin()
	if err != nil {
		return retError("F_00", "", err)
	}
	defer tx.Rollback()

	type tableStruct struct {
		TableName string
		TableVar  any
	}

	var tables []tableStruct
	tables = append(tables,
		tableStruct{"PokemonSpecies", pokemonSpecies},
		tableStruct{"PokemonForm", pokemonForm},
		tableStruct{"PokemonType", pokemonType},
		tableStruct{"PokemonEggGroup", pokemonEggGroup},
	)

	for _, table := range tables {
		err = insertStruct(tx, table.TableName, table.TableVar)
		if err != nil {
			return err
		}
	}

	return nil
}
