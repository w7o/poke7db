package db

/*
Any value stored in here SHOULD BE IMMUTABLE WITHIN THE DATABASE
*/

var mapGrowthClass = map[string]int{
	"medium":      0,
	"erratic":     1,
	"fluctuating": 2,
	"medium-slow": 3,
	"fast":        4,
	"slow":        5,
}

var mapColor = map[string]int{
	"black":  4,
	"blue":   1,
	"brown":  5,
	"gray":   7,
	"green":  3,
	"pink":   9,
	"purple": 6,
	"red":    0,
	"white":  8,
	"yellow": 2,
}

var mapEggGroup = map[int]int{
	1:  0,  // Monster
	8:  1,  // Human-Like
	2:  2,  // Water 1
	9:  3,  // Water 3
	3:  4,  // Bug
	10: 5,  // Mineral
	4:  6,  // Flying
	11: 7,  // Amorphous/Indeterminate
	5:  8,  // Field/Ground
	12: 9,  // Water 2
	6:  10, // Fairy
	13: 11, // Ditto
	7:  12, // Grass/Plant
	14: 13, // Dragon
	15: 14, // No Eggs Discovered
}

var mapShape = map[int]int{
	8:  0,  // Quadruped
	12: 1,  // Humanoid
	6:  2,  // Upright
	2:  3,  // Serpentine
	13: 4,  // Bug-Winged
	9:  5,  // Winged
	14: 6,  // Insectoid
	5:  7,  // Blob
	4:  8,  // Armed
	7:  9,  // Legged
	10: 10, // Multiped
	3:  11, // Fish
	1:  12, // Ball
	11: 13, // Multi-Body
}

var mapType = map[int]int{
	// Values from 1 - 18 are the same value
	19: 127, // Stellar
}
