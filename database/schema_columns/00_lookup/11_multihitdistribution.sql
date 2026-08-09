MultiHitProfileID INTEGER,
HitNumber INTEGER NOT NULL,
HitWeight INTEGER,

PRIMARY KEY (MultiHitProfileID, HitNumber),
FOREIGN KEY (MultiHitProfileID)
    REFERENCES MultiHitProfile(MultiHitProfileID)