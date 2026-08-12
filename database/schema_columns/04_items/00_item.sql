-- @table Item
ItemID INTEGER PRIMARY KEY
    CHECK (ItemID >= 0),
ItemName TEXT NOT NULL,
APIDescription TEXT,
BagCategoryID INTEGER,
FlingEffectID INTEGER,
FlingPower INTEGER,

-- @constraints
FOREIGN KEY (BagCategoryID)
    REFERENCES BagCategory(BagCategoryID),
FOREIGN KEY (FlingEffectID)
    REFERENCES Effect(EffectID)