package main

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"strconv"
)

type PokemonSpeciesDBTable struct {
	PokemonSpeciesID int
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

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Non-recoverable function
func extractIDfromURL(link string) (id int) {
	parsedURL, err := url.Parse(link)
	if err != nil {
		log.Fatal(err)
	}
	id, err = strconv.Atoi(path.Base(parsedURL.Path))
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func speciesToTableStruct(data APIPokemonSpecies) (PokemonSpeciesDBTable, error) {
	colorID, ok := mapColor[string(data.Color)]
	if !ok {
		message := fmt.Errorf("Unknown Pokémon color %s", string(data.Color))
		err := retError("E_00", "", message)
		return PokemonSpeciesDBTable{}, err
	}

	shapeID, ok := mapShape[extractIDfromURL(data.Shape.URL)]
	if !ok {
		message := fmt.Errorf("Unknown Pokémon shape %s with ID %d",
			string(data.Shape.Name), extractIDfromURL(data.Shape.URL))
		err := retError("E_01", "", message)
		return PokemonSpeciesDBTable{}, err
	}

	growthRateID, ok := mapGrowthClass[string(data.GrowthRate)]
	if !ok {
		message := fmt.Errorf("Unknown Pokémon growth rate %s",
			string(data.GrowthRate))
		err := retError("E_02", "", message)
		return PokemonSpeciesDBTable{}, err
	}

	return PokemonSpeciesDBTable{
		PokemonSpeciesID: data.DexNum,
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
