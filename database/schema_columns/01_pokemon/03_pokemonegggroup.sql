-- @table PokemonEggGroup
PokemonSpeciesID INTEGER NOT NULL,
Slot INTEGER NOT NULL
    CHECK (Slot >= 0),
EggGroupID INTEGER NOT NULL,

PRIMARY KEY (PokemonSpeciesID, Slot),
FOREIGN KEY (PokemonSpeciesID)
    REFERENCES PokemonSpecies(PokemonSpeciesID)
FOREIGN KEY (EggGroupID)
    REFERENCES EggGroup(EggGroupID)