package db

type metadataTemplate struct {
	originID   int
	importedAt *string // non-user
	checkedAt  *string // non-user nullable
	createdAt  *string // user
	updatedAt  *string // user nullable
	enabled    int
	hasID      bool // whether or not ID exists
}

type PokemonSpeciesDBEntry struct {
	PokemonSpeciesID string // alphanumeric dex numbers
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

type PokemonEggGroupDBEntry struct {
	PokemonSpeciesID string
	Slot             int
	EggGroupID       int
}

type PokemonEVYieldDBEntry struct {
	PokemonFormID int
	StatID        int
	EVYield       int
}
