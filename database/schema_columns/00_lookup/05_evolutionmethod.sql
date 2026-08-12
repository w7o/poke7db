-- @table EvolutionMethod
EvolutionMethodID INTEGER PRIMARY KEY
    CHECK (EvolutionMethodID >= 0),
Name TEXT NOT NULL, 
Description TEXT