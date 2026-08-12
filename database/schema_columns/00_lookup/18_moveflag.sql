-- @table MoveFlag
MoveFlagID INTEGER PRIMARY KEY
    CHECK (MoveFlagID >= 0),
Name TEXT NOT NULL,
Description TEXT 