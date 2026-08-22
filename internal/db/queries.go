package db

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"strconv"

	"github.com/w7o/poke7db/internal/api"
	"github.com/w7o/poke7db/internal/logging"
)

type tableStruct struct {
	TableName string
	TableVar  any
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

/*
Non-recoverable function
Error codes must be passed through instead of returned
*/
func extractIDfromURL(link string, errCode string) (id int) {
	parsedURL, err := url.Parse(link)
	if err != nil {
		message := fmt.Sprintf("EIFU0 - Failed to parse link %s", link)
		log.Fatal(logging.RetError(errCode, message, err))
	}
	id, err = strconv.Atoi(path.Base(parsedURL.Path))
	if err != nil {
		message := fmt.Sprintf("EIFU1 - Failed to extract integer ID from URL %s",
			parsedURL.String())
		log.Fatal(logging.RetError(errCode, message, err))
	}
	return id
}

func pokemonSpeciesToTS(data api.APIPokemonSpecies) ([]PokemonSpeciesDBEntry, error) {
	colorID, ok := mapColor[string(data.Color)]
	if !ok {
		message := fmt.Errorf("Unknown Pokémon color %s", string(data.Color))
		err := logging.RetError("E_00", "", message)
		return []PokemonSpeciesDBEntry{}, err
	}

	shapeID, ok := mapShape[extractIDfromURL(data.Shape.URL, "E_03")]
	if !ok {
		message := fmt.Errorf("Unknown Pokémon shape %s with ID %d",
			string(data.Shape.Name), extractIDfromURL(data.Shape.URL, ""))
		err := logging.RetError("E_01", "", message)
		return []PokemonSpeciesDBEntry{}, err
	}

	growthRateID, ok := mapGrowthClass[string(data.GrowthRate)]
	if !ok {
		message := fmt.Errorf("Unknown Pokémon growth rate %s",
			string(data.GrowthRate))
		err := logging.RetError("E_02", "", message)
		return []PokemonSpeciesDBEntry{}, err
	}

	return []PokemonSpeciesDBEntry{{
		PokemonSpeciesID: strconv.Itoa(data.DexNum),
		PokemonName:      string(data.Name),
		ColorID:          colorID,
		ShapeID:          shapeID,
		Category:         string(data.Category),
		BaseHappiness:    int(data.BHappiness),
		CaptureRate:      int(data.CaptureRate),
		GrowthRateID:     growthRateID,
		HatchCounter:     data.HatchCounter,
		GenderRate:       int(data.GenderRate) * 2, // frac out of 16 in DB
		IsMythical:       boolToInt(data.IsMythical),
		IsLegendary:      boolToInt(data.IsLegendary),
	}}, nil
}

func pokemonFormsToTS(data api.Pokemon) ([]PokemonFormDBEntry, error) {
	return []PokemonFormDBEntry{{
		PokemonFormID:       data.ID,
		PokemonSpeciesID:    strconv.Itoa(data.DexNum),
		FormName:            data.FormName,
		StatHP:              data.StatBlock.HP.BaseStat,
		StatAttack:          data.StatBlock.Attack.BaseStat,
		StatDefense:         data.StatBlock.Defense.BaseStat,
		StatSpecialAttack:   data.StatBlock.SpecialAttack.BaseStat,
		StatSpecialDefense:  data.StatBlock.SpecialDefense.BaseStat,
		StatSpeed:           data.StatBlock.Speed.BaseStat,
		Height:              data.Height,
		Weight:              data.Weight,
		BaseExperienceYield: data.BaseEXP,
	}}, nil
}

// pass in Pokemon.EggGroups
func pokemonEggGroupToTS(data []api.PokemonEggGroup, pokemonSpeciesID string) ([]PokemonEggGroupDBEntry, error) {
	entry := []PokemonEggGroupDBEntry{}

	for i, eggGroup := range data {
		eggGroupID, ok := mapEggGroup[extractIDfromURL(eggGroup.URL, "E_05")]
		if !ok {
			message := fmt.Errorf("Unknown Pokémon egg group %s with ID %d",
				string(eggGroup.Name), extractIDfromURL(eggGroup.URL, ""))
			err := logging.RetError("E_06", "", message)
			return []PokemonEggGroupDBEntry{}, err
		}
		entry = append(entry, PokemonEggGroupDBEntry{
			PokemonSpeciesID: pokemonSpeciesID,
			Slot:             i,
			EggGroupID:       eggGroupID,
		})
	}
	return entry, nil
}

// pass in Pokemon.StatBlock
func pokemonEVYieldToTS(data api.StatBlock, pokemonFormID int) ([]PokemonEVYieldDBEntry, error) {
	entry := []PokemonEVYieldDBEntry{}

	// helper function
	addEVYield := func(statID int, evYield int) {
		if evYield > 0 {
			entry = append(entry, PokemonEVYieldDBEntry{
				PokemonFormID: pokemonFormID,
				StatID:        statID,
				EVYield:       evYield,
			})
		}
	}

	addEVYield(0, data.HP.EVYield)
	addEVYield(1, data.Attack.EVYield)
	addEVYield(2, data.Defense.EVYield)
	addEVYield(3, data.SpecialAttack.EVYield)
	addEVYield(4, data.SpecialDefense.EVYield)
	addEVYield(5, data.Speed.EVYield)

	return entry, nil
}

func optionalMapI2I(id int, mapping map[int]int) int {
	mappedID, ok := mapping[id]
	if ok {
		return mappedID
	}
	return id
}

// pass in Pokemon.Types
func pokemonTypesToTS(data []api.PokemonType, pokemonFormID int) ([]PokemonTypeDBEntry, error) {
	entry := []PokemonTypeDBEntry{}
	for _, d := range data {
		typeID := extractIDfromURL(d.Type.URL, "E_04")
		typeID = optionalMapI2I(typeID, mapType)

		entry = append(entry,
			PokemonTypeDBEntry{
				PokemonFormID: pokemonFormID,
				Slot:          d.Slot - 1, // 0-indexed in DB
				TypeID:        typeID,
			})
	}

	return entry, nil
}

func DatabasePokemonImport(pokemonData api.Pokemon) error {
	// T01/00
	pokemonSpeciesEntry, err := pokemonSpeciesToTS(*pokemonData.SpeciesInfo)
	if err != nil {
		return err
	}

	// T01/01
	pokemonFormEntry, err := pokemonFormsToTS(pokemonData)
	if err != nil {
		return err
	}

	// T01/02
	typeEntry, err := pokemonTypesToTS(pokemonData.Types, pokemonData.ID)
	if err != nil {
		return err
	}

	// T01/03
	eggGroupEntry, err := pokemonEggGroupToTS(pokemonData.SpeciesInfo.EggGroups,
		strconv.Itoa(pokemonData.DexNum))

	// T01/04
	evYieldEntry, err := pokemonEVYieldToTS(pokemonData.StatBlock, pokemonData.ID)
	if err != nil {
		return err
	}

	// T02/03 (learnset)
	// T03/01 (pokemonability)
	// T04/01 (pokemonhelditem)
	// T10/00 (evolution)

	var tables []tableStruct
	tables = append(tables,
		tableStruct{"PokemonSpecies", pokemonSpeciesEntry},
		tableStruct{"PokemonForm", pokemonFormEntry},
		tableStruct{"PokemonType", typeEntry},
		tableStruct{"PokemonEggGroup", eggGroupEntry},
		tableStruct{"PokemonEVYield", evYieldEntry},
	)

	err = upsertTables(tables, OPokeAPI)
	if err != nil {
		return logging.RetError("E_07",
			fmt.Sprintf("Upserting phase for DatabasePokemonImport for Pokemon No. %d failed",
				pokemonData.DexNum),
			err)
	}

	return nil
}
