package main

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
