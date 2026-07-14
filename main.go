package main

import (
	"fmt"
	"log"
	"os"
)

var SITE string = "https://pokeapi.co/api/v2/"

func main() {
	// if $env:DEV="1" (in powershell) then use locally stored instead
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
}

/* important commands:
go build: creates executable
go get (external link): imports external package and updates go.mod and go.sum
go test -v ./...: test there
*/
