-- @table PokemonForm
PokemonFormID INTEGER PRIMARY KEY 
    CHECK (PokemonFormID >= 0),
PokemonSpeciesID TEXT NOT NULL,
FormName TEXT NOT NULL,	
StatHP INTEGER NOT NULL 
    CHECK (StatHP >= 0),
StatAttack INTEGER NOT NULL 
    CHECK (StatAttack >= 0),
StatDefense INTEGER NOT NULL 
    CHECK (StatDefense >= 0),
StatSpecialAttack INTEGER NOT NULL 
    CHECK (StatSpecialAttack >= 0),
StatSpecialDefense INTEGER NOT NULL 
    CHECK (StatSpecialDefense >= 0),
StatSpeed INTEGER NOT NULL 
    CHECK (StatSpeed >= 0), 
Height INTEGER NOT NULL 
    CHECK (Height >= 0),
Weight INTEGER NOT NULL 
    CHECK (Weight >= 0),
BaseExperienceYield INTEGER NOT NULL 
    CHECK (BaseExperienceYield >= 0),

FOREIGN KEY (PokemonSpeciesID)
    REFERENCES PokemonSpecies(PokemonSpeciesID)