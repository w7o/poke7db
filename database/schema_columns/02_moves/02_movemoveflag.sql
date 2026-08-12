-- @table MoveMoveFlag
MoveID INTEGER NOT NULL,
Slot INTEGER NOT NULL
    CHECK (Slot >= 0),
MoveFlagID INTEGER NOT NULL,

-- @constraints
PRIMARY KEY (MoveID, Slot),
FOREIGN KEY (MoveID)
    REFERENCES Move(MoveID)
FOREIGN KEY (MoveFlagID)
    REFERENCES MoveFlag(MoveFlagID) 