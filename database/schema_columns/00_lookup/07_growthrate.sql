GrowthRateID INTEGER PRIMARY KEY
    CHECK (GrowthRateID >= 0),
Name TEXT NOT NULL,
FormulaLatex TEXT NOT NULL,
Description TEXT