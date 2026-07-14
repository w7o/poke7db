/*
Comments marked with $ are notes for myself in the process of learning Golang
*/

package main

/*
Enum representing growth rate / growth class
Prefix G
*/
type GrowthClass int

const(
	G_MEDIUM_FAST		GrowthClass = iota + 1
	G_ERRATIC
	G_FLUCTUATING
	G_MEDIUM_SLOW
	G_FAST
	G_SLOW
	G_SLIGHT_FAST		GrowthClass = 16
	G_SLIGHT_SLOW		GrowthClass = 17
)


/*
This struct defines a Pokémon's basic statistics derived from their species
*/
type PokemonSpecies struct{
	Name		string		`json:"name"`		// Name of the Pokémon
	FormFlag	int								// Form number of the Pokémon (Default: 0)
	Height		int			`json:"height"`		// Height of the Pokémon, in 0.1kg
	Weight		int			`json:"weight"`		// Weight of the Pokémon, in 0.1kg

	BHappiness	int								// Default happiness, max 255
	CaptureRate	int								// Species-dependent capture probability
	GrowthRate	GrowthClass						// Determines how much XP needed per level
	

	Color		string							// Cosmetic color used within the Pokédex
	Generation	int								// Generation of origin
}