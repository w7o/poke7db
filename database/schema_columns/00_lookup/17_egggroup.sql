-- @table EggGroup
EggGroupID INTEGER PRIMARY KEY
    CHECK (EggGroupID >= 0),
Name TEXT NOT NULL,
Description TEXT 