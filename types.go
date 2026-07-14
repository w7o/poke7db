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
Represents a Pokémon's type
and whether or not it's a maintype (slot = 1) or not
*/
type PokemonType struct {
	Type NamedResource `json:"type"`
	Slot int           `json:"slot"`
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
This string represents the Pokémon's dex color
A separate type is needed for unique unmarshalling logic
*/
type DexColor string

type PokemonEggGroup NamedResource

/*
This struct is the data layout for the /pokemon-species/... endpoint,
that is exported into the main "Pokemon" struct
*/
type APIPokemonSpecies struct {
	BHappiness  int               `json:"base_happiness"`
	CaptureRate int               `json:"capture_rate"`
	Color       DexColor          `json:"color"`
	EggGroups   []PokemonEggGroup `json:"egg_groups"`
}

/*
This struct defines a Pokémon's basic statistics derived from their species
(Both /pokemon and /pokemon-species are merged into this singular struct)
*/
type Pokemon struct {
	Name      string           // Name of the Pokémon species
	DexNum    string           // Pokédex Number, used for calls to pokemon-species endpoint
	FormName  string           `json:"name"`            // Name of the Pokémon form (all lowercase)
	ID        int              `json:"id"`              // PokeAPI ID number
	Height    int              `json:"height"`          // Height of the Pokémon, in 0.1kg
	Weight    int              `json:"weight"`          // Weight of the Pokémon, in 0.1kg
	BaseEXP   int              `json:"base_experience"` // Base experience yield from defeating this Pokémon
	Abilities []PokemonAbility `json:"abilities"`       // Available abilities of a Pokémon (up to 3)
	Moves     []PokemonMove    `json:"moves"`           // Available moves of a Pokémon across all generations
	Types     []PokemonType    `json:"types"`           // The list of types a Pokémon has
	StatBlock StatBlock        `json:"stats"`           // A struct with the six base stats of a species

	BHappiness  int               // Default happiness, max 255
	CaptureRate int               // Species-dependent capture probability
	GrowthRate  GrowthClass       // Determines how much XP needed per level
	EggGroups   []PokemonEggGroup // Which egg group(s) the Pokémon belongs to

	Color      DexColor // Cosmetic color used within the Pokédex
	Generation int      // Generation of origin
	FormFlag   int      // Form number of the Pokémon (Default: 0)

	/*
		Add in future?
		• Held Items
	*/
}
