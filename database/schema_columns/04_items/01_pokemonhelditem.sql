-- @table PokemonHeldItem
PokemonFormID INTEGER NOT NULL,
ItemID INTEGER NOT NULL,
ChanceHeld INTEGER NOT NULL
    CHECK (ChanceHeld BETWEEN 0 AND 100),

PRIMARY KEY (PokemonFormID, ItemID),
FOREIGN KEY (PokemonFormID)
    REFERENCES PokemonForm(PokemonFormID),
FOREIGN KEY (ItemID)
    REFERENCES Item(ItemID)