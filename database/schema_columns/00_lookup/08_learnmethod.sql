-- @table LearnMethod
LearnMethodID INTEGER PRIMARY KEY
    CHECK (LearnMethodID >= 0),
Name TEXT NOT NULL,
Description TEXT