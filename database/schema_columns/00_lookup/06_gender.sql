-- @table Gender
GenderID INTEGER PRIMARY KEY
    CHECK (GenderID >= 0),
Name TEXT NOT NULL