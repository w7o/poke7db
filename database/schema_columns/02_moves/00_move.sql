MoveID INTEGER PRIMARY KEY
    CHECK (MoveID >= 0),
MoveName TEXT NOT NULL,
Accuracy INTEGER NOT NULL
    CHECK (Accuracy >= 0),
CategoryID INTEGER NOT NULL,
PP INTEGER NOT NULL
    CHECK (PP >= 0),
Priority INTEGER NOT NULL,
TypeID INTEGER NOT NULL,
Power INTEGER NOT NULL
    CHECK (Power >= 0),
DamageClassID INTEGER NOT NULL,
TargetID INTEGER NOT NULL
    CHECK (Target >= 0),
CritBonus INTEGER NOT NULL
    CHECK (CritBonus >= 0),
MultiHitTypeID INTEGER NOT NULL,
MoveDescription TEXT NOT NULL

FOREIGN KEY (CategoryID)
    REFERENCES MoveCategory(MoveCategoryID),
FOREIGN KEY (TypeID)
    REFERENCES Type(TypeID),
FOREIGN KEY (DamageClassID)
    REFERENCES DamageClass(DamageClassID),
FOREIGN KEY (TargetID)
    REFERENCES Target(TargetID),
FOREIGN KEY (MultiHitTypeID)
    REFERENCES MultiHitProfile(MultiHitProfileID)