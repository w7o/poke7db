-- @table Stat
StatID INTEGER PRIMARY KEY
    CHECK (StatID >= 0),
Name TEXT NOT NULL,
Description TEXT 