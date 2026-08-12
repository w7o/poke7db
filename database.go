package main

import (
	// "context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

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
	ID INTEGER UNIQUE
		CHECK (ID IS NOT NULL OR enabled = FALSE),

	FOREIGN KEY (originID)
		REFERENCES DataOrigin(DataOriginID)
`

// Initialize Database as a database
var database *sql.DB

/*
Note: Need colMetadata = "" if no metadata
*/
func generateTableQuery(tableName, colBody, colMetadata string) string {
	if colMetadata != "" {
		colMetadata = ", " + colMetadata
	}
	return fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		%s %s
	)
	`, tableName, colBody, colMetadata)
}

func grabDBFlag(sqlFile string, flag string) (string, error) {
	flag = "-- @" + flag // e.g. -- @table

	for _, line := range strings.Split(sqlFile, "\n") {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, flag) {
			// removes the flag then removes any whitespace around it
			stringReturn := strings.TrimSpace(strings.TrimPrefix(line, flag))
			return stringReturn, nil
		}
	}

	return "", fmt.Errorf("Missing %s metadata", flag)
}

func grabTableName(sqlFile string) (string, error) {
	return grabDBFlag(sqlFile, "table")
}

func databasePing(database *sql.DB) error {
	// Test that the database works
	err := database.Ping()
	if err != nil {
		return retError("D_06", "Database ping failed", err)
	}
	return nil
}

/*
Index 00 is hardcoded to not have metadata applied
*/
func databaseCreateLookup(database *sql.DB) error {
	// Lookup tables have the same metadata as User Main tables
	log.Print("DB: Creating lookup tables")
	err := databasePing(database)
	if err != nil {
		return err
	}

	tx, err := database.Begin()
	if err != nil {
		return retError("D_00", "", err)
	}
	defer tx.Rollback()

	// Create the rest of the lookup tables
	var dir string = "./database/schema_columns/00_lookup"
	// Golang ReadDir already sorts in alphabetical order
	files, err := os.ReadDir(dir)
	if err != nil {
		return retError("D_01", "", err)
	}

	// regex object to grab filename index numbers
	var re_fileNameIndex *regexp.Regexp = regexp.MustCompile("^([0-9]{2})_.+$")

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		matches := re_fileNameIndex.FindStringSubmatch(file.Name())

		if len(matches) < 2 {
			log.Printf("DB: Notice: Ignoring SQL file in lookup %s with incorrect filename format", file.Name())
			continue
		}

		// convert obtained index to int
		s, err := strconv.Atoi(matches[1])
		if err != nil {
			return retError("D_02", "Failed to convert an index to int", err)
		}

		// obtain the sql file
		schemaColumns, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return retError("D_03", "Failed to obtain .sql file", err)
		}

		// obtain table name from sql file
		tableName, err := grabTableName(string(schemaColumns))
		if err != nil {
			return retError("D_04", "Failed to obtain table name from .sql file", err)
		}

		query := generateTableQuery(tableName, string(schemaColumns), userMetadata)
		// No metadata for T00/00 DataOrigin specifically
		if s == 0 {
			query = generateTableQuery(tableName, string(schemaColumns), "")
		}
		_, err = tx.Exec(query)
		if err != nil {
			writeError(query)
			message := fmt.Sprintf("Failed to create table %s from %s; \n\toffending query outputted to error.txt\n\t",
				tableName, file.Name())
			return retError("D_05", message, err)
		}
	}
	return tx.Commit()
}

/*
Initializes or opens the database.
If initializing:
• Creates lookup tables and fills them with required information (from database/seed)
• Creates main tables and user tables
*/
func DatabaseInit(filepath string) (*sql.DB, error) {
	log.Print("DB: Initializing SQLite database")

	var err error

	// Create database if the database at the filepath doesn't exist
	database, err = sql.Open("sqlite", filepath)
	if err != nil {
		return nil, err
	}

	err = databasePing(database)
	if err != nil {
		return nil, err
	}

	// Turn on foreign keys
	_, err = database.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}

	// Creates lookup tables and fills them with required information
	err = databaseCreateLookup(database)
	if err := databaseCreateLookup(database); err != nil {
		return nil, err
	}

	fmt.Println("SQLite database initialized")

	return database, nil
}
