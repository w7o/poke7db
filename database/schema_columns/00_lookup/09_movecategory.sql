MoveCategoryID INTEGER PRIMARY KEY
    CHECK (MoveCategoryID >= 0),
Name TEXT NOT NULL,
Description TEXT