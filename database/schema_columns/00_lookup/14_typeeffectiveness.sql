-- @table TypeEffectiveness
-- @noID
TypeATK INTEGER,
TypeDEF INTEGER,
Multiplier INTEGER NOT NULL,

-- @constraints
PRIMARY KEY (TypeATK, TypeDEF),
FOREIGN KEY (TypeATK)
    REFERENCES Type(TypeID),
FOREIGN KEY (TypeDEF)
    REFERENCES Type(TypeID)