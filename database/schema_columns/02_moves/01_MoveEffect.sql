-- @table MoveEffect
-- @noID
MoveID INTEGER NOT NULL,
EffectOrder INTEGER NOT NULL
    CHECK (EffectOrder >= 1),
EffectID INTEGER NOT NULL,
EffectChance INTEGER NOT NULL
    CHECK (EffectChance BETWEEN 0 AND 100),
EffectValue INTEGER NOT NULL,

-- @constraints
PRIMARY KEY (MoveID, EffectOrder),
FOREIGN KEY (MoveID)
    REFERENCES Move(MoveID),
FOREIGN KEY (EffectID)
    REFERENCES Effect(EffectID)