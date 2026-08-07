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
	{"MoveEffects", moveEffectsTemplate},			// TU007
	{"PokemonTypes", pokemonTypesTemplate},			// TU008
	{"Evolutions", evolutionsTemplate},				// TU009
}

const endingNonUser string = `
	originID INTEGER NOT NULL,
	importedAt TEXT NOT NULL,
	checkedAt TEXT,
	enabled INTEGER NOT NULL DEFAULT 1 
		CHECK (enabled IN (0, 1))

	FOREIGN KEY (originID)
		REFERENCES DataOrigin(DataOriginID)
`

const endingUser string = `
	originID INTEGER NOT NULL,
	createdAt TEXT NOT NULL,
	updatedAt TEXT,
	enabled INTEGER NOT NULL DEFAULT 1 
		CHECK (enabled IN (0, 1))

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
	GrowthRate INTEGER NOT NULL
		CHECK (GrowthRate >= 0),
	GenderRate INTEGER NOT NULL
		CHECK (GrowthRate >= -1),
	HatchCounter INTEGER NOT NULL
        CHECK (HatchCounter >= 0),
    Color INTEGER NOT NULL,
    Shape INTEGER NOT NULL,
    IsMythical INTEGER NOT NULL
        CHECK (IsMythical IN (0, 1)),
    IsLegendary INTEGER NOT NULL
        CHECK (IsLegendary IN (0, 1))
	
	FOREIGN KEY (Color)
		REFERENCES Color(ColorID)
	FOREIGN KEY (Shape)
		REFERENCES Shape(ShapeID)
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
		CHECK (BaseExperienceYield >= 0)

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
    AbliityName TEXT NOT NULL,
    EffectDescription TEXT NOT NULL
`

const learnsetTemplate string = `
    PokemonFormID INTEGER NOT NULL,
		CHECK (PokemonFormID >= 0)
    MoveID INTEGER NOT NULL,
		CHECK (MoveID >= 0),
    LearnMethod INTEGER NOT NULL
        CHECK (LearnMethod >= 0),
    Level INTEGER NOT NULL
        CHECK (Level >= 0),
    LearnDescription TEXT NOT NULL,

    PRIMARY KEY (PokemonFormID, MoveID, LearnMethod, Level),
    FOREIGN KEY (PokemonFormID)
        REFERENCES PokemonForm(PokemonFormID),
    FOREIGN KEY (MoveID)
        REFERENCES Move(MoveID)
	FOREIGN KEY (LearnMethod)
		REFERENCES LearnMethod(LearnMethodID)
`

const moveTemplate string = `
    MoveID INTEGER PRIMARY KEY
        CHECK (MoveID >= 0),
    MoveName TEXT NOT NULL,
    Accuracy INTEGER NOT NULL
        CHECK (Accuracy >= 0),
    Category INTEGER NOT NULL,
    PP INTEGER NOT NULL
        CHECK (PP >= 0),
    Priority INTEGER NOT NULL,
    Type INTEGER NOT NULL,
    Power INTEGER NOT NULL
        CHECK (Power >= 0),
    DamageClass INTEGER NOT NULL,
    Target INTEGER NOT NULL
        CHECK (Target >= 0),
    CritBonus INTEGER NOT NULL
        CHECK (CritBonus >= 0),
    MultiHitType INTEGER NOT NULL,
    MoveDescription TEXT NOT NULL

	FOREIGN KEY (Category)
		REFERENCES MoveCategory(MoveCategoryID),
	FOREIGN KEY (Type)
		REFERENCES Type(TypeID),
	FOREIGN KEY (DamageClass)
		REFERENCES DamageClass(DamageClassID),
	FOREIGN KEY (Target)
		REFERENCES Target(TargetID),
	FOREIGN KEY (MultiHitType)
		REFERENCES MultiHitProfile(MultiHitProfileID)
`

const moveEffectsTemplate = `
    MoveID INTEGER NOT NULL,
    EffectOrder INTEGER NOT NULL
        CHECK (EffectOrder >= 1),
    Effect INTEGER NOT NULL,
    EffectChance INTEGER NOT NULL
        CHECK (EffectChance BETWEEN 0 AND 100),
    EffectValue INTEGER NOT NULL,

    PRIMARY KEY (MoveID, EffectOrder),
    FOREIGN KEY (MoveID)
        REFERENCES Move(MoveID),
    FOREIGN KEY (Effect)
        REFERENCES Effect(EffectID)
`

const pokemonTypesTemplate = `
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

const evolutionsTemplate = `
	EvolutionID INTEGER PRIMARY KEY,

    PokemonSpeciesFrom TEXT NOT NULL,
    PokemonSpeciesTo TEXT NOT NULL,
    Method INTEGER NOT NULL,
    Level INTEGER NOT NULL
        CHECK (Level >= 0),
    Gender INTEGER,
    ItemID INTEGER,
    TimeOfDay INTEGER,
    MinimumValue INTEGER,
    Location TEXT,
    KnownMoveID INTEGER,
    KnownMoveTypeID INTEGER,

    FOREIGN KEY (PokemonSpeciesFrom)
        REFERENCES PokemonSpecies(PokemonSpeciesID),
    FOREIGN KEY (PokemonSpeciesTo)
        REFERENCES PokemonSpecies(PokemonSpeciesID),
	FOREIGN KEY (Gender)
		REFERENCES Gender(GenderID),
	FOREIGN KEY (Method)
		REFERENCES EvolutionMethod(EvolutionMethodID),
    FOREIGN KEY (ItemID)
        REFERENCES Item(ItemID),
	FOREIGN KEY (TimeOfDay)
		REFERENCES TimeOfDay(TimeID),
    FOREIGN KEY (KnownMoveID)
        REFERENCES Move(MoveID),
    FOREIGN KEY (KnownMoveTypeID)
        REFERENCES Type(TypeID)
`

const pokemonHeldItemTemplate = `
	PokemonFormID INTEGER,
	ItemID INTEGER,
	ChanceHeld INTEGER
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
	ItemName TEXT,
	APIDescription TEXT,
	BagCategory INTEGER, 
	FlingEffect INTEGER,
	FlingPower INTEGER,
	
	FOREIGN KEY (BagCategory)
		REFERENCES BagCategory(BagCategoryID),
	FOREIGN KEY (FlingEffect)
		REFERENCES Effect(EffectID)
`

const pokemonEggGroupTemplate =`
	PokemonSpeciesID INTEGER,
	Slot INTEGER
		CHECK (Slot >= 0),
	EggGroupID INTEGER,
	
	PRIMARY KEY (PokemonSpeciesID, Slot),
	FOREIGN KEY (PokemonSpeciesID)
		REFERENCES PokemonSpecies(PokemonSpeciesID)
	FOREIGN KEY (EggGroupID)
		REFERENCES EggGroup(EggGroupID)
`

// TODO: Design and add MoveFlag lookup table
const moveMoveFlagTemplate = `
	MoveID INTEGER,
	Slot INTEGER
		CHECK (Slot >= 0),
	MoveFlagID INTEGER,

	PRIMARY KEY (MoveID, Slot),
	FOREIGN KEY (MoveFlagID)
		REFERENCES MoveFlag(MoveFlagID) -- not yet implemented
`