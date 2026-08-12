/*
Usage:
make run
	Runs the program against the API using the default Pokémon (ID 197 / Umbreon)

make run ID=<pokemon>
Runs the program against the API using the specified Pokémon ID or name
	Example:
		make run ID=ampharos

make server
	Starts the local API instance. Must be run in a separate terminal, and must
	be run before run-dev.

make run-dev
	Runs the program against a local API instance using the default Pokémon

make run-dev ID=<pokemon>
	Runs the program against a local API instance using the specified Pokémon
	Example:
		make run-dev ID=espeon
*/

package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/davecgh/go-spew/spew"
)

var IS_DEV_BUILD bool = true // MUST BE SET TO FALSE FOR EVERY VERSION BUMP
var SITE string = "https://pokeapi.co/api/v2/"
var VERSION string = "0.0.9"
var PROJECT_NAME string = "Poké7DB"

func generateVersionNumber() {
	// If not development build, don't print commit details
	// %TODO pre-release / beta/ alpha / release tag support
	if !IS_DEV_BUILD {
		VERSION = fmt.Sprintf("%s v%s", PROJECT_NAME, VERSION)
		return
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		VERSION = fmt.Sprintf("%s UNKNOWN VERSION", PROJECT_NAME)
		return
	}

	commit := "unknown"
	modified := ""
	revTime := ""

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.time":
			t, err := time.Parse(time.RFC3339, setting.Value)
			if err != nil {
				fmt.Println(err)
				return
			}
			revTime = fmt.Sprintf("%d", t.Unix())
		case "vcs.revision":
			commit = setting.Value[:7]
		case "vcs.modified":
			modified = " (modified)"
		}
	}

	VERSION = fmt.Sprintf("%s v%s.dev-%s_%s%s", PROJECT_NAME, VERSION, revTime, commit, modified)
}

func main() {
	clearErrorFile()

	generateVersionNumber()
	fmt.Println(VERSION)

	// if $env:P7D_ENV="dev" then use locally stored instead

	if os.Getenv("P7D_ENV") == "dev" {
		SITE = "http://localhost:8080/"
	}

	log.SetPrefix("pk7db: ")
	// Format: Ldate, Ltime, Lshortfile
	log.SetFlags(0b11001000)

	// testing variable (defaults to 197 as per Makefile)
	pokemonID := os.Args[1]

	// initialize database
	database, err := DatabaseInit("./database/app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	data, err := poke_query(pokemonID)

	if err != nil {
		log.Fatal(err.Error())
	}

	// fmt.Printf("%+v\n", data)
	spew.Dump(data)

	fmt.Printf("\nFinished request on %s\n\n%s\n", SITE, VERSION)
}

/* important commands:
go build: creates executable
go get (external link): imports external package and updates go.mod and go.sum
go test -v ./...: test there
*/
