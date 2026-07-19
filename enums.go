package main

/*
Enum representing growth rate / growth class
Prefix G
*/
type GrowthClass int

const (
	G_MEDIUM_FAST GrowthClass = iota + 1
	G_ERRATIC
	G_FLUCTUATING
	G_MEDIUM_SLOW
	G_FAST
	G_SLOW

	// Unused growth rates from Gen 1 & 2 games
	G_SLIGHT_FAST GrowthClass = 16
	G_SLIGHT_SLOW GrowthClass = 17
)
