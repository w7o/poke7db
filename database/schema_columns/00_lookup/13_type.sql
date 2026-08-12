-- @table Type
TypeID INTEGER PRIMARY KEY
    CHECK (TypeID >= 0),
Name TEXT NOT NULL,
Description TEXT 