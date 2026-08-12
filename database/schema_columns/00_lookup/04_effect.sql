-- @table Effect
EffectID INTEGER PRIMARY KEY
    CHECK (EffectID >= 0),
Name TEXT NOT NULL,
Description TEXT