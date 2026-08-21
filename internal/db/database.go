package db

import (
	// "context"

	"bytes"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/w7o/poke7db/internal/logging"

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

const userMetadataMain string = `
	originID INTEGER NOT NULL,
	createdAt TEXT NOT NULL,
	updatedAt TEXT,
	enabled INTEGER NOT NULL DEFAULT 1 
		CHECK (enabled IN (0, 1)),
`
const userMetadataMainConstraints string = `
	FOREIGN KEY (originID)
		REFERENCES DataOrigin(DataOriginID)
`

const userMetadataID string = `
	ID INTEGER UNIQUE
		CHECK (ID IS NOT NULL OR enabled = FALSE),
`

const userMetadata string = userMetadataMain + userMetadataMainConstraints
const userMetadataWithID string = userMetadataMain + userMetadataID + userMetadataMainConstraints

// Initialize Database as a database
var database *sql.DB

// Convert []string to []any
func convStringsToAnys(values []string) []any {
	var result []any
	for _, value := range values {
		result = append(result, value)
	}
	return result
}

// Convert sql.Rows to []string
func convArbitraryRowsToStrings(rows *sql.Rows, numCol int) ([]string, error) {
	var strings []string

	for rows.Next() {
		// since Scan requires pointers,
		// create fixed length lists to support variable pointers
		values := make([]any, numCol)
		pointers := make([]any, numCol)
		for i := range values {
			pointers[i] = &values[i]
		}

		// scans rows into the pointers, which themselves point to a table of values
		err := rows.Scan(pointers...)
		if err != nil {
			return nil, err
		}

		// doing type assertion instead of make([]string) in case i want
		// to split that into another helper function later
		for _, value := range values {
			str, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("Expected primary key name to be string, got %T", value)
			}
			if value != nil {
				strings = append(strings, str)
			}
		}

	}
	return strings, nil
}

func extractFileIndexNumber(filename string) (int, error) {
	var reFileNameIndex *regexp.Regexp = regexp.MustCompile("^([0-9]{2})_.+$")

	matches := reFileNameIndex.FindStringSubmatch(filename)
	// matches[0] = entire matched STRING | matches[1] = capture group ([0-9]{2})
	if len(matches) < 2 {
		log.Printf("DB: Notice: Ignoring SQL file in lookup %s with incorrect filename format", filename)
		return -1, nil
	}

	// convert obtained index to int
	s, err := strconv.Atoi(matches[1])
	if err != nil {
		return -1, err
	}
	return s, nil
}

func grabDBFlagParameter(sqlFile string, flag string) (string, error) {
	for _, line := range strings.Split(sqlFile, "\n") {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, flag) {
			// removes the flag then removes any whitespace around it
			stringReturn := strings.TrimSpace(strings.TrimPrefix(line, flag))
			return stringReturn, nil
		}
	}

	return "", fmt.Errorf("D_19: Missing %s metadata", flag)
}

func fileFormatToFlag(flagText string, fileFormat string) string {
	switch fileFormat {
	case "sql":
		return "-- @" + flagText
	case "csv": // explicit
		return "# @" + flagText
	default:
		return "# @" + flagText
	}
}

// -- @table <tableName>
func grabTableName(fileContents string, fileFormat string) (string, error) {
	flag := fileFormatToFlag("table", fileFormat)
	return grabDBFlagParameter(fileContents, flag)
}

// -- @noID
func checkIncludeID(fileContents string, fileFormat string) (bool, error) {
	flag := fileFormatToFlag("noID", fileFormat)
	_, err := grabDBFlagParameter(fileContents, flag)
	if err != nil {
		// D_19 reports that there's no flag found which is allowable
		// Basically catching D_19 'exception'
		if strings.HasPrefix(err.Error(), "D_19") {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func grabConstraintsSQL(sqlFileContents string) (body string, constraints string) {
	flag := "-- @constraints"

	lines := strings.Split(sqlFileContents, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, flag) {
			body := strings.TrimSpace(strings.Join(lines[:i], "\n"))
			body = strings.TrimSuffix(body, ",") // remove trailing commas

			rest := strings.TrimSpace(strings.TrimPrefix(line, flag))
			if i+1 < len(lines) {
				rest += "\n" + strings.Join(lines[i+1:], "\n")
			}

			return body, strings.TrimSpace(rest)
		}
	}

	return strings.TrimSpace(sqlFileContents), ""
}

/*
Note: Need colMetadata = "" if no metadata
*/
func generateTableQuery(tableName, colBody, colMetadata string) (string, error) {
	colBody, colConstraint := grabConstraintsSQL(colBody)

	// Add commas as necessary
	if colMetadata != "" || colConstraint != "" {
		colBody += ","
	}
	if colMetadata != "" && colConstraint != "" {
		colMetadata += ","
	}

	return fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		%s %s %s
	)
	`, tableName, colBody, colMetadata, colConstraint), nil
}

func stringRows(rows *sql.Rows) (string, error) {
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}

	// the shit i need to do to just print out a query dude
	var buffer bytes.Buffer

	t := table.NewWriter()
	t.SetOutputMirror(&buffer)
	var header table.Row
	for _, column := range columns {
		header = append(header, column)
	}

	// Table header
	t.AppendHeader(header)

	for rows.Next() {
		// since Scan requires pointers, create fixed length lists to support variable pointers
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		// scans rows into the pointers, which themselves point to a table of values
		err := rows.Scan(pointers...)
		if err != nil {
			return "", err
		}

		var row table.Row
		for _, value := range values {
			if value == nil {
				row = append(row, "NULL")
			} else {
				row = append(row, value)
			}
		}
		t.AppendRow(row)
	}

	err = rows.Err()
	if err != nil {
		return "", err
	}

	t.Render()

	return buffer.String(), nil
}

func databasePing(database *sql.DB) error {
	// Test that the database works
	err := database.Ping()
	if err != nil {
		return logging.RetError("D_06", "Database ping failed", err)
	}
	return nil
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
	if err != nil {
		return nil, err
	}
	logging.WriteLog("DB: Created lookup tables")

	err = databaseSeedLookup(database)
	if err != nil {
		return nil, err
	}
	logging.WriteLog("DB: Seeded lookup tables")

	err = databaseCreateMain(database)
	if err != nil {
		return nil, err
	}
	logging.WriteLog("DB: Created main tables")

	fmt.Println("SQLite database initialized")

	return database, nil
}

/*
temporary query
*/
func SampleDatabaseQuery(database *sql.DB) error {
	rows, err := database.Query("SELECT * FROM MoveFlag")
	if err != nil {
		return err
	}

	defer rows.Close()
	ret, err := stringRows(rows)
	if err != nil {
		return err
	}

	logging.ResetMessageFile("logDB", "logDB.txt")
	logging.WriteFile(ret, "logDB.txt")
	return nil
}
