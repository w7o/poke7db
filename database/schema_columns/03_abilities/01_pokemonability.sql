-- @table PokemonAbility
PokemonFormID INTEGER NOT NULL,
Slot INTEGER NOT NULL
    CHECK (Slot >= 0),
IsHidden INTEGER NOT NULL
    CHECK (IsHidden IN (0, 1)),
AbilityID INTEGER NOT NULL,

PRIMARY KEY (PokemonFormID, Slot),
FOREIGN KEY (PokemonFormID)
    REFERENCES PokemonForm(PokemonFormID),
FOREIGN KEY (AbilityID)
    REFERENCES Ability(AbilityID)