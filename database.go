package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// Initialize Database as a database
var Database *sql.DB

func DatabaseInit(filepath string){
	fmt.Println("Initializing SQLite database")
	
	Database, err := sql.Open("sqlite", filepath)
	if err != nil {
		log.Fatal(err)
	}

	err = Database.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("SQLite database initialized")
}