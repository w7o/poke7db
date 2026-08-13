package main

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"strconv"
)

type PokemonSpeciesDBEntry struct {
	PokemonSpeciesID string
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

type PokemonFormDBEntry struct {
	PokemonFormID       int
	FormName            string
	PokemonSpeciesID    string
	Height              int
	Weight              int
	BaseExperienceYield int
	StatHP              int
	StatAttack          int
	StatDefense         int
	StatSpecialAttack   int
	StatSpecialDefense  int
	StatSpeed           int
}

type PokemonTypeDBEntry struct {
	PokemonFormID int
	Slot          int
	TypeID        int
}

type PokemonEVYieldDBEntry struct {
	PokemonFormID int
	StatID        int
	EVYield       int
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
		log.Fatal(retError(errCode, message, err))
	}
	id, err = strconv.Atoi(path.Base(parsedURL.Path))
	if err != nil {
		message := fmt.Sprintf("EIFU1 - Failed to extract integer ID from URL %s",
			parsedURL.String())
		log.Fatal(retError(errCode, message, err))
	}
	return id
}

func pokemonSpeciesToTS(data APIPokemonSpecies) (PokemonSpeciesDBEntry, error) {
	colorID, ok := mapColor[string(data.Color)]
	if !ok {
		message := fmt.Errorf("Unknown Pokémon color %s", string(data.Color))
		err := retError("E_00", "", message)
		return PokemonSpeciesDBEntry{}, err
	}

	shapeID, ok := mapShape[extractIDfromURL(data.Shape.URL, "E_03")]
	if !ok {
		message := fmt.Errorf("Unknown Pokémon shape %s with ID %d",
			string(data.Shape.Name), extractIDfromURL(data.Shape.URL, ""))
		err := retError("E_01", "", message)
		return PokemonSpeciesDBEntry{}, err
	}

	growthRateID, ok := mapGrowthClass[string(data.GrowthRate)]
	if !ok {
		message := fmt.Errorf("Unknown Pokémon growth rate %s",
			string(data.GrowthRate))
		err := retError("E_02", "", message)
		return PokemonSpeciesDBEntry{}, err
	}

	return PokemonSpeciesDBEntry{
		PokemonSpeciesID: strconv.Itoa(data.DexNum),
		PokemonName:      string(data.Name),
		ColorID:          colorID,
		ShapeID:          shapeID,
		Category:         string(data.Category),
		BaseHappiness:    int(data.BHappiness),
		CaptureRate:      int(data.CaptureRate),
		GrowthRateID:     growthRateID,
		HatchCounter:     data.HatchCounter,
		IsMythical:       boolToInt(data.IsMythical),
		IsLegendary:      boolToInt(data.IsLegendary),
	}, nil
}

func pokemonFormsToTS(data Pokemon) (PokemonFormDBEntry, error) {
	return PokemonFormDBEntry{
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
	}, nil
}

// pass in Pokemon.StatBlock
func pokemonEVYieldToTS(data StatBlock, pokemonFormID int) ([]PokemonEVYieldDBEntry, error) {
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

// pass in Pokemon.Types
func pokemonTypesToTS(data []PokemonType, pokemonFormID int) ([]PokemonTypeDBEntry, error) {
	entry := []PokemonTypeDBEntry{}
	for _, d := range data {
		entry = append(entry,
			PokemonTypeDBEntry{
				PokemonFormID: pokemonFormID,
				Slot:          d.Slot,
				TypeID:        extractIDfromURL(d.Type.URL, "E_04") + 1,
				// pokeapi is 1-indexed, db is 0-indexed
			})
	}

	return entry, nil
}

// note: doesn't actually import to database yet ඞ
func DatabasePokemonImport(pokemonData Pokemon) ([]any, error) {
	// T01/00
	pokemonSpeciesEntry, err := pokemonSpeciesToTS(*pokemonData.SpeciesInfo)
	if err != nil {
		return nil, err
	}

	// T01/01
	pokemonFormEntry, err := pokemonFormsToTS(pokemonData)
	if err != nil {
		return nil, err
	}

	// T01/02
	typeEntry, err := pokemonTypesToTS(pokemonData.Types, pokemonData.ID)
	if err != nil {
		return nil, err
	}

	// T01/03 @TODO
	// eggGroupEntry, err := pokemonEggGroupToTS(pokemonData.SpeciesInfo.EggGroups)

	// T01/04
	evYieldEntry, err := pokemonEVYieldToTS(pokemonData.StatBlock, pokemonData.ID)
	if err != nil {
		return nil, err
	}

	// T02/03 (learnset)
	// T03/01 (pokemonability)
	// T04/01 (pokemonhelditem)
	// T10/00 (evolution)

	ret := []any{}
	ret = append(ret,
		pokemonSpeciesEntry,
		pokemonFormEntry,
		typeEntry,
		evYieldEntry)
	return ret, nil
}
