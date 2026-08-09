/*
Refer to ./design/...
	- CS_719-08.pdf for Main Table outlines
	- CS_719-07.pdf for Reference Table outlines
*/

package main

// type TableDefinition struct {
// 	Name 	string 	// Name of the table
// 	Template string	// Major columns of the table
// }

// var majorTables = []TableDefinition{
// 	{"PokemonForm", pokemonFormTemplate},			// TU001
// 	{"PokemonSpecies", pokemonSpeciesTemplate},		// TU002
// 	{"PokemonAbility", pokemonAbilityTemplate},		// TU003
// 	{"Ability", abilityTemplate},					// TU004
// 	{"Learnset", learnsetTemplate}, 				// TU005
// 	{"Move", moveTemplate},							// TU006
// 	{"MoveEffect", moveEffectTemplate},				// TU007
// 	{"PokemonType", pokemonTypeTemplate},			// TU008
// 	{"Evolution", evolutionTemplate},				// TU009
// 	{"PokemonHeldItem", pokemonHeldItemTemplate},	// TU010
// 	{"Item", itemTemplate},							// TU011
// 	{"PokemonEggGroup", pokemonEggGroupTemplate},	// TU012
// 	{"MoveMoveFlag", moveMoveFlagTemplate},			// TU013
// }

const nonUserMetadata string = `
	originID INTEGER NOT NULL,
	importedAt TEXT NOT NULL,
	checkedAt TEXT,
	enabled INTEGER NOT NULL DEFAULT 1 
		CHECK (enabled IN (0, 1)),

	FOREIGN KEY (originID)
		REFERENCES DataOrigin(DataOriginID)
`

const userMetadata string = `
	originID INTEGER NOT NULL,
	createdAt TEXT NOT NULL,
	updatedAt TEXT,
	enabled INTEGER NOT NULL DEFAULT 1 
		CHECK (enabled IN (0, 1)),

	FOREIGN KEY (originID)
		REFERENCES DataOrigin(DataOriginID)
`