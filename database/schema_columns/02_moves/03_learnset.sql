-- @table Learnset
PokemonFormID INTEGER NOT NULL,
MoveID INTEGER NOT NULL,
LearnMethodID INTEGER NOT NULL,
Level INTEGER
    CHECK (Level >= 0),
LearnDescription TEXT NOT NULL,

PRIMARY KEY (PokemonFormID, MoveID, LearnMethod, Level),
FOREIGN KEY (PokemonFormID)
    REFERENCES PokemonForm(PokemonFormID),
FOREIGN KEY (MoveID)
    REFERENCES Move(MoveID)
FOREIGN KEY (LearnMethodID)
    REFERENCES LearnMethod(LearnMethodID)