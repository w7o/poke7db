-- @table PokemonSpecies
PokemonSpeciesID TEXT PRIMARY KEY,
PokemonName TEXT NOT NULL,
Category TEXT NOT NULL,
BaseHappiness INTEGER NOT NULL
    CHECK (BaseHappiness >= 0),
CaptureRate	INTEGER NOT NULL
    CHECK (CaptureRate >= 0),
GrowthRateID INTEGER NOT NULL
    CHECK (GrowthRate >= 0),
GenderRate INTEGER NOT NULL
    CHECK (GenderRate >= -1),
HatchCounter INTEGER NOT NULL
    CHECK (HatchCounter >= 0),
ColorID INTEGER NOT NULL,
ShapeID INTEGER NOT NULL,
IsMythical INTEGER NOT NULL
    CHECK (IsMythical IN (0, 1)),
IsLegendary INTEGER NOT NULL
    CHECK (IsLegendary IN (0, 1)),

FOREIGN KEY (ColorID)
    REFERENCES Color(ColorID),
FOREIGN KEY (ShapeID)
    REFERENCES Shape(ShapeID),
FOREIGN KEY (GrowthRateID)
    REFERENCES GrowthRate(GrowthRateID)