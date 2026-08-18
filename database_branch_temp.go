package main

import (
	"database/sql"
	"fmt"
	"reflect"
	// "strconv"
	"strings"
	"time"
)

var pokemonSpecies []PokemonSpeciesDBEntry
var pokemonForm []PokemonFormDBEntry
var pokemonType []PokemonTypeDBEntry
var pokemonEggGroup []PokemonEggGroupDBEntry
var pokemonEVYields []PokemonEVYieldDBEntry

func nullableString(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func insertStruct(tx *sql.Tx, tableName string, dbStruct any,
	metadata metadataTemplate) error {

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

	// set up table column names
	for i := range fields {
		fieldNames = append(fieldNames, structType.Field(i).Name)
	}

	// set up metadata
	// sanity check
	hasNonUser := metadata.importedAt != nil || metadata.checkedAt != nil
	hasUser := metadata.createdAt != nil || metadata.updatedAt != nil
	hasUserCheck := hasUser || metadata.hasID
	if hasNonUser == hasUserCheck {
		return retError("F_05", "Invalid metadata field combination", nil)
	}

	var metadataValues []any
	fieldNames = append(fieldNames, "originID", "enabled")
	metadataValues = append(metadataValues,
		metadata.originID, //originID
		1,                 //enabled
	)
	if hasNonUser {
		fieldNames = append(fieldNames, "importedAt", "checkedAt")
		metadataValues = append(metadataValues,
			nullableString(metadata.importedAt), //importedAt
			nullableString(metadata.checkedAt),  //checkedAt
		)
	} else { // hasUser
		fieldNames = append(fieldNames, "createdAt", "updatedAt")
		metadataValues = append(metadataValues,
			nullableString(metadata.createdAt), //createdAt
			nullableString(metadata.updatedAt), //updatedAt
		)
		if metadata.hasID {
			fieldNames = append(fieldNames, "ID")
			// metadata value insertion done in below loop
		}
	}

	writeLog(fmt.Sprintf("DB: Field names: %v", fieldNames))
	var placeholderIndex int = 1
	for i := 0; i < structValue.Len(); i++ {
		// currStruct represents the current struct (i.e. table)
		currStruct := structValue.Index(i)
		var rowPlaceholders []string

		// add placeholders for the current row
		for j := range fields {
			fieldValues = append(fieldValues, currStruct.Field(j).Interface())
			rowPlaceholders = append(rowPlaceholders,
				fmt.Sprintf("$%d", placeholderIndex))
			placeholderIndex++
		}

		// add metadata values
		for _, value := range metadataValues {
			fieldValues = append(fieldValues, value)
			rowPlaceholders = append(rowPlaceholders,
				fmt.Sprintf("$%d", placeholderIndex))
			placeholderIndex++
		}

		// add ID w/ correct value if it's enabled
		if metadata.hasID {
			fieldValues = append(
				fieldValues,
				currStruct.Field(0).Interface(),
			)
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
	writeLog(fmt.Sprintf("DB:\n\tquery:%s\n\tvalues: %v", query, fieldValues))
	if err != nil {
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

	timestamp := time.Now().UTC().Format(time.RFC3339)
	metadata := metadataTemplate{
		originID:   1, // PokéAPI
		importedAt: &timestamp,
		checkedAt:  nil,
		enabled:    1,
		hasID:      false,
	}

	var tables []tableStruct
	tables = append(tables,
		tableStruct{"PokemonSpecies", pokemonSpecies},
		tableStruct{"PokemonForm", pokemonForm},
		tableStruct{"PokemonType", pokemonType},
		tableStruct{"PokemonEggGroup", pokemonEggGroup},
		tableStruct{"PokemonEVYield", pokemonEVYields},
	)

	for _, table := range tables {
		err = insertStruct(tx, table.TableName, table.TableVar, metadata)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
