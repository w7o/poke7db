package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"
)

var DEV bool = true // used for version number
var SITE string = "https://pokeapi.co/api/v2/"
var VERSION string = "0.0.1"
var PROJECTNAME string = "Poké7DB"

func generateVersionNumber() {
	// If not development version, don't print commit details
	if !DEV {
		VERSION = fmt.Sprintf("%s v%s", PROJECTNAME, VERSION)
		return
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		VERSION = fmt.Sprintf("%s UNKNOWN VERSION", PROJECTNAME)
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

	VERSION = fmt.Sprintf("%s v%s-%s-%s%s", PROJECTNAME, VERSION, revTime, commit, modified)
}

func main() {
	// if $env:DEV="1" (in powershell) then use locally stored instead

	generateVersionNumber()
	fmt.Println(VERSION)

	if os.Getenv("DEV") == "1" {
		SITE = "http://localhost:8080/"
	}

	log.SetPrefix("pk7db: ")
	//Format: Ldate, Ltime, Lshortfile
	log.SetFlags(0b11001000)

	// testing variable
	pokemonID := "197"

	data, err := poke_query(pokemonID)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("%+v\n", data)

	fmt.Printf("\n%s\n", VERSION)
}

/* important commands:
go build: creates executable
go get (external link): imports external package and updates go.mod and go.sum
go test -v ./...: test there
*/
