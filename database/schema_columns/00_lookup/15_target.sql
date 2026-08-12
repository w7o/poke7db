-- @table Target
TargetID INTEGER PRIMARY KEY
    CHECK (TargetID >= 0),
Name TEXT NOT NULL,
Description TEXT 