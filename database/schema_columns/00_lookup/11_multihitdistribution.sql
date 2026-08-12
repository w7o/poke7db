-- @table MultiHitDistribution
-- @noID
MultiHitProfileID INTEGER,
HitNumber INTEGER NOT NULL,
HitWeight INTEGER,

-- @constraints
PRIMARY KEY (MultiHitProfileID, HitNumber),
FOREIGN KEY (MultiHitProfileID)
    REFERENCES MultiHitProfile(MultiHitProfileID)
