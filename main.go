package main

import (
	"fmt"
	"log"
)


func main() {
	log.SetPrefix("pk7db: ")
	//Format: Ldate, Ltime, Lshortfile
	log.SetFlags(0b11001000)

	pokemon := "197"

	data, err := poke_query(pokemon)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Print(data)
}

/* important commands:
go build: creates executable
go get (external link): imports external package and updates go.mod and go.sum
go test -v ./...: test there 
*/
