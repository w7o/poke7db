/*
Refer to ./design/...
	- CS_719-08.pdf for Main Table outlines
	- CS_719-07.pdf for Reference Table outlines
*/

package main

type TableDefinition struct {
	Name 	string 	// Name of the table
	Template string	// Major columns of the table
}

var majorTables = []TableDefinition{
	{"PokemonForm", pokemonFormTemplate},			// TU001
	{"PokemonSpecies", pokemonSpeciesTemplate},		// TU002
	{"PokemonAbility", pokemonAbilityTemplate},		// TU003
	{"Ability", abilityTemplate},					// TU004
	{"Learnset", learnsetTemplate}, 				// TU005
	{"Move", moveTemplate},							// TU006
	{"MoveEffect", moveEffectTemplate},				// TU007
	{"PokemonType", pokemonTypeTemplate},			// TU008
	{"Evolution", evolutionTemplate},				// TU009
	{"PokemonHeldItem", pokemonHeldItemTemplate},	// TU010
	{"Item", itemTemplate},							// TU011
	{"PokemonEggGroup", pokemonEggGroupTemplate},	// TU012
	{"MoveMoveFlag", moveMoveFlagTemplate},			// TU013
}

const nonUserMetadata string = `
	originID INTEGER NOT NULL,
	importedAt TEXT NOT NULL,
	checkedAt TEXT,
	enabled INTEGER NOT NULL DEFAULT 1 
		CHECK (enabled IN (0, 1)),

	FOREIGN KEY (originID)
		REFERENCES DataOrigin(DataOriginID)
`

const userMetadata string = `
	originID INTEGER NOT NULL,
	createdAt TEXT NOT NULL,
	updatedAt TEXT,
	enabled INTEGER NOT NULL DEFAULT 1 
		CHECK (enabled IN (0, 1)),

	FOREIGN KEY (originID)
		REFERENCES DataOrigin(DataOriginID)
`

const pokemonSpeciesTemplate string = `
	PokemonSpeciesID TEXT PRIMARY KEY,
	PokemonName TEXT NOT NULL,
	Category TEXT NOT NULL,
	BaseHappiness INTEGER NOT NULL
		CHECK (BaseHappiness >= 0),
	CaptureRate	INTEGER NOT NULL
		CHECK (CaptureRate >= 0),
	GrowthRateID INTEGER NOT NULL
		CHECK (GrowthRate >= 0),
	GenderRate INTEGER NOT NULL
		CHECK (GenderRate >= -1),
	HatchCounter INTEGER NOT NULL
        CHECK (HatchCounter >= 0),
    ColorID INTEGER NOT NULL,
    ShapeID INTEGER NOT NULL,
    IsMythical INTEGER NOT NULL
        CHECK (IsMythical IN (0, 1)),
    IsLegendary INTEGER NOT NULL
        CHECK (IsLegendary IN (0, 1)),
	
	FOREIGN KEY (ColorID)
		REFERENCES Color(ColorID)
	FOREIGN KEY (ShapeID)
		REFERENCES Shape(ShapeID)
	FOREIGN KEY (GrowthRateID)
		REFERENCES GrowthRate(GrowthRateID)
`

const pokemonFormTemplate string = `
	PokemonFormID INTEGER PRIMARY KEY 
		CHECK (PokemonFormID >= 0),
	PokemonSpeciesID TEXT NOT NULL,
	FormName TEXT NOT NULL,	
	StatHP INTEGER NOT NULL 
		CHECK (StatHP >= 0),
	StatAttack INTEGER NOT NULL 
		CHECK (StatAttack >= 0),
	StatDefense INTEGER NOT NULL 
		CHECK (StatDefense >= 0),
	StatSpecialAttack INTEGER NOT NULL 
		CHECK (StatSpecialAttack >= 0),
	StatSpecialDefense INTEGER NOT NULL 
		CHECK (StatSpecialDefense >= 0),
	StatSpeed INTEGER NOT NULL 
		CHECK (StatSpeed >= 0), 
	Height INTEGER NOT NULL 
		CHECK (Height >= 0),
	Weight INTEGER NOT NULL 
		CHECK (Weight >= 0),
	BaseExperienceYield INTEGER NOT NULL 
		CHECK (BaseExperienceYield >= 0),

	FOREIGN KEY (PokemonSpeciesID)
		REFERENCES PokemonSpecies(PokemonSpeciesID)
` 

const pokemonAbilityTemplate string = `
    PokemonFormID INTEGER NOT NULL,
    Slot INTEGER NOT NULL
        CHECK (Slot >= 0),
    IsHidden INTEGER NOT NULL
        CHECK (IsHidden IN (0, 1)),
    AbilityID INTEGER NOT NULL,

    PRIMARY KEY (PokemonFormID, Slot),
    FOREIGN KEY (PokemonFormID)
        REFERENCES PokemonForm(PokemonFormID),
    FOREIGN KEY (AbilityID)
        REFERENCES Ability(AbilityID)
`

const abilityTemplate string = `
    AbilityID INTEGER PRIMARY KEY
        CHECK (AbilityID >= 0),
    AbilityName TEXT NOT NULL,
    EffectDescription TEXT NOT NULL
`

const learnsetTemplate string = `
    PokemonFormID INTEGER NOT NULL,
    MoveID INTEGER NOT NULL,
    LearnMethodID INTEGER NOT NULL,
    Level INTEGER
        CHECK (Level >= 0),
    LearnDescription TEXT NOT NULL,

    PRIMARY KEY (PokemonFormID, MoveID, LearnMethod, Level),
    FOREIGN KEY (PokemonFormID)
        REFERENCES PokemonForm(PokemonFormID),
    FOREIGN KEY (MoveID)
        REFERENCES Move(MoveID)
	FOREIGN KEY (LearnMethodID)
		REFERENCES LearnMethod(LearnMethodID)
`

const moveTemplate string = `
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
`

const moveEffectTemplate = `
    MoveID INTEGER NOT NULL,
    EffectOrder INTEGER NOT NULL
        CHECK (EffectOrder >= 1),
    EffectID INTEGER NOT NULL,
    EffectChance INTEGER NOT NULL
        CHECK (EffectChance BETWEEN 0 AND 100),
    EffectValue INTEGER NOT NULL,

    PRIMARY KEY (MoveID, EffectOrder),
    FOREIGN KEY (MoveID)
        REFERENCES Move(MoveID),
    FOREIGN KEY (EffectID)
        REFERENCES Effect(EffectID)
`

const pokemonTypeTemplate = `
    PokemonFormID INTEGER NOT NULL,
    Slot INTEGER NOT NULL
        CHECK (Slot >= 0),
    TypeID INTEGER NOT NULL,

    PRIMARY KEY (PokemonFormID, Slot),
    FOREIGN KEY (PokemonFormID)
        REFERENCES PokemonForm(PokemonFormID),
    FOREIGN KEY (TypeID)
        REFERENCES Type(TypeID)
`

const evolutionTemplate = `
	EvolutionID INTEGER PRIMARY KEY,

    PokemonSpeciesFROM TEXT NOT NULL,
    PokemonSpeciesTO TEXT NOT NULL,
    MethodID INTEGER NOT NULL,
    Level INTEGER NOT NULL
        CHECK (Level >= 0),
    GenderID INTEGER,
    ItemID INTEGER,
    TimeOfDayID INTEGER,
    MinimumValue INTEGER,
    Location TEXT,
    KnownMoveID INTEGER,
    KnownMoveTypeID INTEGER,

    FOREIGN KEY (PokemonSpeciesFROM)
        REFERENCES PokemonSpecies(PokemonSpeciesID),
    FOREIGN KEY (PokemonSpeciesTO)
        REFERENCES PokemonSpecies(PokemonSpeciesID),
	FOREIGN KEY (GenderID)
		REFERENCES Gender(GenderID),
	FOREIGN KEY (MethodID)
		REFERENCES EvolutionMethod(EvolutionMethodID),
    FOREIGN KEY (ItemID)
        REFERENCES Item(ItemID),
	FOREIGN KEY (TimeOfDayID)
		REFERENCES TimeOfDay(TimeID),
    FOREIGN KEY (KnownMoveID)
        REFERENCES Move(MoveID),
    FOREIGN KEY (KnownMoveTypeID)
        REFERENCES Type(TypeID)
`

const pokemonHeldItemTemplate = `
	PokemonFormID INTEGER NOT NULL,
	ItemID INTEGER NOT NULL,
	ChanceHeld INTEGER NOT NULL
		CHECK (ChanceHeld BETWEEN 0 AND 100),

	PRIMARY KEY (PokemonFormID, ItemID),
	FOREIGN KEY (PokemonFormID)
		REFERENCES PokemonForm(PokemonFormID),
	FOREIGN KEY (ItemID)
		REFERENCES Item(ItemID)
`

const itemTemplate =`
	ItemID INTEGER PRIMARY KEY
		CHECK (ItemID >= 0),
	ItemName TEXT NOT NULL,
	APIDescription TEXT,
	BagCategoryID INTEGER,
	FlingEffectID INTEGER,
	FlingPower INTEGER,
	
	FOREIGN KEY (BagCategoryID)
		REFERENCES BagCategory(BagCategoryID),
	FOREIGN KEY (FlingEffectID)
		REFERENCES Effect(EffectID)
`

const pokemonEggGroupTemplate =`
	PokemonSpeciesID INTEGER NOT NULL,
	Slot INTEGER NOT NULL
		CHECK (Slot >= 0),
	EggGroupID INTEGER NOT NULL,
	
	PRIMARY KEY (PokemonSpeciesID, Slot),
	FOREIGN KEY (PokemonSpeciesID)
		REFERENCES PokemonSpecies(PokemonSpeciesID)
	FOREIGN KEY (EggGroupID)
		REFERENCES EggGroup(EggGroupID)
`

// TODO: Design and add MoveFlag lookup table
const moveMoveFlagTemplate = `
	MoveID INTEGER NOT NULL,
	Slot INTEGER NOT NULL
		CHECK (Slot >= 0),
	MoveFlagID INTEGER NOT NULL,

	PRIMARY KEY (MoveID, Slot),
	FOREIGN KEY (MoveID)
		REFERENCES Move(MoveID)
	FOREIGN KEY (MoveFlagID)
		REFERENCES MoveFlag(MoveFlagID) -- not yet implemented
`