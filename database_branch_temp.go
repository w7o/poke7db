package main

/*
Temporary branch because I accidentally left my v0.1.4 - v0.1.5 commits unpushed
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

func initTemporaryData() {
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

}
