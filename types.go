package main

/*
Enum representing growth rate / growth class
Prefix G
*/
type GrowthClass int

const (
	G_MEDIUM_FAST GrowthClass = iota + 1
	G_ERRATIC
	G_FLUCTUATING
	G_MEDIUM_SLOW
	G_FAST
	G_SLOW
	G_SLIGHT_FAST GrowthClass = 16
	G_SLIGHT_SLOW GrowthClass = 17
)

/*
Represents a generic object with a name and URL to that object's information
*/
type NamedResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

/*
Represents an ability's species-specific
properties:
• Whether or not it's a hidden ability
• Which slot the ability is in (is this even important)
*/
type PokemonAbility struct {
	Ability  NamedResource `json:"ability"`
	IsHidden bool          `json:"is_hidden"`
	Slot     int           `json:"slot"`
}

/*
Represents a move's species-specific
properties, though none exist at the time of implementation
*/
type PokemonMove struct {
	Move NamedResource `json:"move"`
}

/*
Represents a Pokémon's base stat as well as the EV yield from
beating that Pokémon
*/
type PokemonStat struct {
	BaseStat int
	EVYield  int
}

/*
This struct contains the six base stats of a Pokémon species
*/
type StatBlock struct {
	HP             PokemonStat
	Attack         PokemonStat
	Defense        PokemonStat
	SpecialAttack  PokemonStat
	SpecialDefense PokemonStat
	Speed          PokemonStat
}

/*
This struct defines a Pokémon's basic statistics derived from their species
*/
type PokemonSpecies struct {
	Name      string           `json:"name"`            // Name of the Pokémon
	Height    int              `json:"height"`          // Height of the Pokémon, in 0.1kg
	Weight    int              `json:"weight"`          // Weight of the Pokémon, in 0.1kg
	BaseEXP   int              `json:"base_experience"` // Base experience yield from defeating this Pokémon
	Abilities []PokemonAbility `json:"abilities"`       // Available abilities of a Pokémon (up to 3)
	Moves     []PokemonMove    `json:"moves"`           // Available moves of a Pokémon across all generations
	StatBlock StatBlock        `json:"stats"`           // A struct with the six base stats of a species

	FormFlag int // Form number of the Pokémon (Default: 0)

	BHappiness  int         // Default happiness, max 255
	CaptureRate int         // Species-dependent capture probability
	GrowthRate  GrowthClass // Determines how much XP needed per level

	Color      string // Cosmetic color used within the Pokédex
	Generation int    // Generation of origin

	/*
		Add in future?
		• Held Items
	*/
}
