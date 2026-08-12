-- @table TypeEffectiveness
TypeATK INTEGER,
TypeDEF INTEGER,
Multiplier INTEGER NOT NULL,

-- @constraints
PRIMARY KEY (TypeATK, TypeDEF),
FOREIGN KEY (TypeATK)
    REFERENCES Type(TypeATK),
FOREIGN KEY (TypeDEF)
    REFERENCES Type(TypeDEF)