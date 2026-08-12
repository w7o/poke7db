-- @table Evolution
EvolutionID INTEGER PRIMARY KEY,

PokemonSpeciesFROM TEXT NOT NULL,
PokemonSpeciesTO TEXT NOT NULL,
MethodID INTEGER NOT NULL,
Level INTEGER NOT NULL
    CHECK (Level >= 0),
GenderID INTEGER,
ItemID INTEGER,
TimeOfDayID INTEGER,
MinimumValue INTEGER,
Location TEXT,
KnownMoveID INTEGER,
KnownMoveTypeID INTEGER,

FOREIGN KEY (PokemonSpeciesFROM)
    REFERENCES PokemonSpecies(PokemonSpeciesID),
FOREIGN KEY (PokemonSpeciesTO)
    REFERENCES PokemonSpecies(PokemonSpeciesID),
FOREIGN KEY (GenderID)
    REFERENCES Gender(GenderID),
FOREIGN KEY (MethodID)
    REFERENCES EvolutionMethod(EvolutionMethodID),
FOREIGN KEY (ItemID)
    REFERENCES Item(ItemID),
FOREIGN KEY (TimeOfDayID)
    REFERENCES TimeOfDay(TimeID),
FOREIGN KEY (KnownMoveID)
    REFERENCES Move(MoveID),
FOREIGN KEY (KnownMoveTypeID)
    REFERENCES Type(TypeID)