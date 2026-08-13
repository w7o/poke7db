-- @table PokemonEVYield
-- @noID
PokemonFormID INTEGER NOT NULL,
StatID INTEGER NOT NULL,
EVYield INTEGER NOT NULL
    CHECK (EVYield >= 0),

-- @constraints
PRIMARY KEY (PokemonFormID, StatID),
FOREIGN KEY (PokemonFormID)
    REFERENCES PokemonForm(PokemonFormID),
FOREIGN KEY (StatID)
    REFERENCES Stat(StatID)