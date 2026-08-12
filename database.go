package main

import (
	// "context"
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

const nonUserMetadata string = `
	originID INTEGER NOT NULL,
	importedAt TEXT NOT NULL,
	checkedAt TEXT,
	enabled INTEGER NOT NULL DEFAULT 1 
		CHECK (enabled IN (0, 1)),

	FOREIGN KEY (originID)
		REFERENCES DataOrigin(DataOriginID)
`

const userMetadata string = `
	originID INTEGER NOT NULL,
	createdAt TEXT NOT NULL,
	updatedAt TEXT,
	enabled INTEGER NOT NULL DEFAULT 1 
		CHECK (enabled IN (0, 1)),
	ID INTEGER UNIQUE,
		CHECK (ID IS NOT NULL OR enabled = FALSE)

	FOREIGN KEY (originID)
		REFERENCES DataOrigin(DataOriginID)
`

// Initialize Database as a database
var database *sql.DB

func createTableQuery(tableName, template, colMetadata string) string {
	return fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
	%s
	,
	%s
	)
	`, tableName, template, colMetadata)
}

func DatabaseInit(filepath string){
	fmt.Println("Initializing SQLite database")
	
	var err error

	database, err = sql.Open("sqlite", filepath)
	if err != nil {
		log.Fatal(err)
	}

	err = database.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Turn on foreign keys
	_, err = database.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}

	// var query string = createTableQuery(
	// 	"PokemonForm", pokemonFormTemplate, nonUserMetadata)

	// // Create PokemonForm database
	// _, err = database.ExecContext(
	// 	context.Background(),
	// 	query,
	// )
	// if err != nil{
	// 	log.Fatal(err)
	// }
	fmt.Println("SQLite database initialized")
}