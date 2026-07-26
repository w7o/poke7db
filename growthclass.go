package main

type GrowthClass string

/*
import (
	"fmt"
	"strings"
)
*/

/*
Enum representing growth rate / growth class
Prefix G
*/

//go:generate stringer -type=GrowthClass
/*
type GrowthClass int

const (
	G_MEDIUM_FAST GrowthClass = iota
	G_ERRATIC
	G_FLUCTUATING
	G_MEDIUM_SLOW
	G_FAST
	G_SLOW

	// Unused growth rates from Gen 1 & 2 games
	G_SLIGHT_FAST GrowthClass = 16
	G_SLIGHT_SLOW GrowthClass = 17
)

var growthClassAPIString = map[string]GrowthClass{
	"medium":              G_MEDIUM_FAST,
	"slow-then-very-fast": G_ERRATIC,
	"fast-then-very-slow": G_FLUCTUATING,
	"medium-slow":         G_MEDIUM_SLOW,
	"fast":                G_FAST,
	"slow":                G_SLOW,
}

func ParseAPIGrowthClass(apiString string) (GrowthClass, error) {
	growth, ok := growthClassAPIString[strings.ToLower(apiString)]
	if !ok {
		return 0, fmt.Errorf("Unknown growth class \"%s\"", apiString)
	}
	return growth, nil
}
*/
