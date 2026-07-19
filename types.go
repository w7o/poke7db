package main

/*
Various types and structs
*/

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

type PokemonForm struct {
	Name      string
	APINumber int
}

/*
This string represents the Pokémon's dex color
A separate type is needed for unique unmarshalling logic
*/

type Category string
type DexColor string
type EvolutionChain NamedResource
type Generation int
type Language NamedResource
type PokemonEggGroup NamedResource
type PokemonName string
type PokemonShape NamedResource

/*
This struct represents the Pokémon's next evolutions
*/
type PokemonEvolutions struct {
	URL        string          // Will allow the choosing of evolution method later
	Evolutions []NamedResource // List of Pokémon with their own species URLs
	// ^ Note: Evolution methods may vary between main series games
}

/*
This struct is the data layout for the /pokemon-species/... endpoint,
that is exported into the main "Pokemon" struct
*/
type APIPokemonSpecies struct {
	BHappiness   uint8             `json:"base_happiness"` // 0 - 255
	CaptureRate  uint8             `json:"capture_rate"`   // 0 - 255
	Color        DexColor          `json:"color"`
	EggGroups    []PokemonEggGroup `json:"egg_groups"`
	GenderRate   int8              `json:"gender_rate"` // -1 - 8
	Category     Category          `json:"genera"`
	Generation   Generation        `json:"generation"`
	HatchCounter int               `json:"hatch_counter"`
	IsMythical   bool              `json:"is_mythical"`
	IsLegendary  bool              `json:"is_legendary"`
	Name         PokemonName       `json:"names"`
	Shape        PokemonShape      `json:"shape"`

	// for further processing
	EvolutionChain EvolutionChain `json:"evolution_chain"`
	// FormsList      []PokemonForm  `json:"varieties"`
}

/*
This struct defines a Pokémon's basic statistics derived from their species
(Both /pokemon and /pokemon-species are merged into this singular struct)
*/
type Pokemon struct {
	FormName  string           `json:"name"`            // Name of the Pokémon form (all lowercase)
	ID        int              `json:"id"`              // PokeAPI ID number (different per form)
	Height    int              `json:"height"`          // Height of the Pokémon, in 0.1kg
	Weight    int              `json:"weight"`          // Weight of the Pokémon, in 0.1kg
	BaseEXP   int              `json:"base_experience"` // Base experience yield from defeating this Pokémon
	Abilities []PokemonAbility `json:"abilities"`       // Available abilities of a Pokémon (up to 3)
	Moves     []PokemonMove    `json:"moves"`           // Available moves of a Pokémon across all generations
	Types     []PokemonType    `json:"types"`           // The list of types a Pokémon has
	StatBlock StatBlock        `json:"stats"`           // A struct with the six base stats of a species

	// Species-dependent, i.e. doesn't depend on a Pokémon's form
	Name         PokemonName         // Name of the Pokémon species in proper capitalization
	DexNum       string              // Pokédex Number, used for calls to pokemon-species endpoint
	BHappiness   uint8               // Default happiness, max 255
	CaptureRate  uint8               // Species-dependent base capture probability, max 255
	Category     Category            // A label for the Pokémon, e.g. the "Moonlight Pokémon" for Umbreon
	GenderRate   int8                // Female gender probability, in 1/8ths (e.g. 3 means 3/8ths chance ♀)
	GrowthRate   GrowthClass         // Determines how much XP needed per level
	EggGroups    []PokemonEggGroup   // Which egg group(s) the Pokémon belongs to
	Evolutions   []PokemonEvolutions // Pokémon the current Pokémon can evolve to
	HatchCounter int                 // Base factor for determining species-dependent steps to hatch egg
	Shape        PokemonShape        // The shape of a Pokémon as defined in some Pokédexes
	IsMythical   bool                // Is the Pokémon a mythical
	IsLegendary  bool                // Is the Pokémon a legendary
	// %TODO Pokemon Evolutions Parsing
	// %TODO Growth rate

	Color      DexColor   // Cosmetic color used within the Pokédex
	Generation Generation // Generation of origin

	// Dex Entries handled and stored at the database side

	/*
		Add in future?
		• Held Items
	*/
}
