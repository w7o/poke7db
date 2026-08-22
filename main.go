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
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/w7o/poke7db/internal/api"
	"github.com/w7o/poke7db/internal/db"
	"github.com/w7o/poke7db/internal/logging"
	"github.com/w7o/poke7db/internal/version"
)

func main() {
	version.GenerateVersionNumber()
	fmt.Println(version.G_Version)

	logging.ResetLogFile()

	// if $env:P7D_ENV="dev" then use locally stored instead

	if os.Getenv("P7D_ENV") == "dev" {
		version.G_Site = "http://localhost:8080/"
	}

	log.SetPrefix("pk7db: ")
	// Format: Ldate, Ltime, Lshortfile
	log.SetFlags(0b11001000)

	// testing variable (defaults to 197 as per Makefile)
	pokemonID := os.Args[1]
	db.D_setDexNum(pokemonID) // debug print

	// initialize database
	database, err := db.DatabaseInit("./database/app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	var pokemonData api.Pokemon
	pokemonData, err = api.Poke_Query(pokemonID)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("\nFinished request on %s", version.G_Site)

	db.SampleDatabaseQuery(database)
	fmt.Printf("DATABASE EXAMPLE OUTPUT TO logDB.txt\n")

	var dataString bytes.Buffer
	spew.Fdump(&dataString, pokemonData)

	logging.ResetMessageFile("Other logs", "logOther.txt")
	logging.WriteFile(dataString.String(), "logOther.txt")

	fmt.Printf("API OUTPUT TO logOther.txt\n")
	if os.Getenv("P7D_WRITE_DB") == "0" {
		logging.WriteLog("p7d_write_db no")
		os.Exit(0)
	}
	logging.WriteLog("p7d_write_db yes")

	err = db.DatabasePokemonImport(pokemonData)
	if err != nil {
		log.Fatal(err)
	}

	// logging.ResetMessageFile("logQuery", "logQuery.txt")
	// for _, item := range info {
	// 	spew.Fdump(&dataString, item)
	// 	logging.WriteFile(dataString.String(), "logQuery.txt")
	// }

	// fmt.Printf("TEST OUTPUT TO logQuery.txt\n")

	// err = db.InitTemporaryData()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	fmt.Printf("\n\n%s", version.G_Version)
	os.Exit(0)
}

/* important commands:
go build: creates executable
go get (external link): imports external package and updates go.mod and go.sum
go test -v ./...: test there
*/
