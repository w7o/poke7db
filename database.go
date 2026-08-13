package main

import (
	// "context"

	"bytes"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

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
func strToAny(values []string) []any {
	var result []any
	for _, value := range values {
		result = append(result, value)
	}
	return result
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

func seedTable(tx *sql.Tx, csvPath string, sqlPath string, tableIndex int, timestamp string) error {
	// Condition to exclude metadata
	includeMetadataCondition := tableIndex != 0

	csvContents, err := os.ReadFile(csvPath)
	if err != nil {
		return retError("D_11", "Failed to obtain .csv file", err)
	}

	sqlContents, err := os.ReadFile(sqlPath)
	if err != nil {
		message := fmt.Sprintf("csv -> sql file \"%s\" doesn't appear to exist", sqlPath)
		return retError("D_21", message, err)
	}

	tableName, err := grabTableName(string(csvContents), "csv")
	if err != nil {
		return retError("D_12", "Failed to grab table name csv", err)
	}

	includeID, err := checkIncludeID(string(sqlContents), "sql")
	if err != nil {
		return retError("D_20", "Failed to grab ID inclusion bool sql", err)
	}

	csvFile, _ := os.Open(csvPath)

	// Create a new reader for the CSV file
	reader := csv.NewReader(csvFile)
	reader.Comment = '#'

	// CSV Header
	columns, err := reader.Read()
	if err != nil {
		return retError("D_13", "Failed to read CSV header", err)
	}

	// Include createdAt timestamp if table contains metadata
	// NOTE: If editing to include more seeded metadata in the future, check
	// all includeMetadataCondition usages
	const userMetadataSeedValuesCount int = 2

	if includeMetadataCondition {
		columns = append(
			append([]string{}, columns...),
			"createdAt",
			"enabled",
		)
		if includeID {
			columns = append(
				append([]string{}, columns...),
				"ID",
			)
		}
	}

	var placeholders []string
	for range columns {
		placeholders = append(placeholders, "?")
	}

	// Ready a placeholder query
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) ON CONFLICT DO NOTHING",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)
	writeLog(fmt.Sprintf("DB: Query -  %s", query))
	statement, err := tx.Prepare(query)
	if err != nil {
		return retError("D_15", "Placeholder query preparation failed", err)
	}

	defer statement.Close()

	// repeat until break (end of file)
	for {
		row, err := reader.Read()

		// end of file
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return retError("D_14", "Failed to read CSV row", err)
		}

		// row/column count sanity check
		expectedValuesCount := len(columns)
		if includeMetadataCondition {
			expectedValuesCount -= userMetadataSeedValuesCount
			if includeID {
				expectedValuesCount -= 1
			}
		}
		if len(row) != expectedValuesCount {
			message := fmt.Sprintf("CSV row has %d values, expected %d", len(row), expectedValuesCount)
			return retError("D_16", message, err)
		}

		args := strToAny(row)

		if includeMetadataCondition {
			// ..., createdAt, enabled
			args = append(args, timestamp, 1)
			if includeID {
				args = append(args, args[0])
			}
		}

		writeLog(fmt.Sprintf("DB: Arguments - %v", args))

		r, err := statement.Exec(args...)
		if err != nil {
			return retError("D_17", "Executing statement failed", err)
		}

		rowsAffected, err := r.RowsAffected()
		if err != nil {
			return retError("D_18", "Failed to get rows affected", err)
		}

		writeLog(fmt.Sprintf("DB: %d row[s] affected", rowsAffected))
	}
	return nil
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
		return retError("D_06", "Database ping failed", err)
	}
	return nil
}

/*
Index 00 is hardcoded to not have metadata applied
*/
func databaseCreateLookup(database *sql.DB) error {
	// Lookup tables have the same metadata as User Main tables
	writeLogAndConsole("DB: Creating lookup tables")
	err := databasePing(database)
	if err != nil {
		return err
	}

	tx, err := database.Begin()
	if err != nil {
		return retError("D_00", "", err)
	}
	defer tx.Rollback()

	// Create the lookup tables
	var dir string = "./database/schema_columns/00_lookup"
	// Golang ReadDir already sorts in alphabetical order
	files, err := os.ReadDir(dir)
	if err != nil {
		return retError("D_01", "", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		s, err := extractFileIndexNumber(file.Name())
		if err != nil {
			return retError("D_02", "Failed to extract filename index number", err)
		}
		if s == -1 {
			continue
		}

		// obtain the sql file
		schemaColumns, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return retError("D_03", "Failed to obtain .sql file", err)
		}

		// obtain table name from sql file
		tableName, err := grabTableName(string(schemaColumns), "sql")
		if err != nil {
			return retError("D_04", "Failed to obtain table name from .sql file", err)
		}

		// obtain include ID bool from sql file
		includeID, err := checkIncludeID(string(schemaColumns), "sql")
		if err != nil {
			return retError("D_22", "Failed to obtain include ID bool from .sql file", err)
		}

		// set No Metadata Condition
		var noMetadataCondition bool = s == 0

		// set correct metadata value
		var metadata string = userMetadata
		if noMetadataCondition {
			metadata = ""
		} else if includeID {
			metadata = userMetadataWithID
		}

		query, err := generateTableQuery(tableName, string(schemaColumns), metadata)
		// No metadata for T00/00 DataOrigin specifically
		if err != nil {
			return retError("D_07", "Error when splitting body and constraint", err)
		}

		_, err = tx.Exec(query)
		if err != nil {
			writeLog(query)
			message := fmt.Sprintf("Failed to create table %s from %s; \n\toffending query outputted to error.txt\n\t",
				tableName, file.Name())
			return retError("D_05", message, err)
		}
	}
	return tx.Commit()
}

func clearLookupEntries(tx *sql.Tx, files []os.DirEntry) error {
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) != ".csv" {
			continue
		}

		// check if valid file
		s, err := extractFileIndexNumber(file.Name())
		if err != nil {
			return retError("D_26", "Failed to extract file index number", err)
		}
		if s == -1 {
			continue
		}

		var seedDir string = "./database/seed/00_lookup"
		csvPath := filepath.Join(seedDir, file.Name())
		csv, err := os.ReadFile(csvPath)
		if err != nil {
			return retError("D_35", "Failed to read CSV file", err)
		}

		tableName, err := grabTableName(string(csv), "csv")
		if err != nil {
			return retError("D_36", "Failed to read CSV table name", err)
		}

		clearQuery := fmt.Sprintf("DELETE FROM %s", tableName)
		_, err = tx.Exec(clearQuery)
		if err != nil {
			return retError("D_37", fmt.Sprintf("Failed to clear table %s",
				tableName), err)
		}
	}
	return nil
}

func databaseSeedLookup(database *sql.DB) error {
	writeLogAndConsole("DB: Seeding lookup tables")
	err := databasePing(database)
	if err != nil {
		return err
	}

	tx, err := database.Begin()
	if err != nil {
		return retError("D_08", "Database transaction init failed", err)
	}
	defer tx.Rollback()

	var seedDir string = "./database/seed/00_lookup"
	var schemaDir string = "./database/schema_columns/00_lookup"

	files, err := os.ReadDir(seedDir)
	if err != nil {
		return retError("D_09", "Reading seed lookup directory failed", err)
	}

	//@USER make this a setting/optional at some point
	err = clearLookupEntries(tx, files)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) != ".csv" {
			continue
		}
		sqlFileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())) + ".sql"

		s, err := extractFileIndexNumber(file.Name())
		if err != nil {
			return retError("D_10", "Failed to convert an index to int", err)
		}
		if s == -1 {
			continue
		}

		// obtain the csv file
		csvPath := filepath.Join(seedDir, file.Name())

		sqlPath := filepath.Join(schemaDir, sqlFileName)

		timestamp := time.Now().UTC().Format(time.RFC3339)

		err = seedTable(tx, csvPath, sqlPath, s, timestamp)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func databaseCreateMain(database *sql.DB) error {
	writeLogAndConsole("DB: Creating main tables")
	err := databasePing(database)
	if err != nil {
		return err
	}

	tx, err := database.Begin()
	if err != nil {
		return retError("D_23", "Transaction start fail", err)
	}
	defer tx.Rollback()

	var schemaDir string = "./database/schema_columns"
	var schemaUserDir string = "./database/schema_columns_user"

	files, err := os.ReadDir(schemaDir)
	if err != nil {
		return retError("D_24", "Read schema col dir failed", err)
	}

	var dirNames []string
	for _, dir := range files {
		if dir.IsDir() {
			if dir.Name() != "00_lookup" {
				dirNames = append(dirNames, dir.Name())
			}
		}
	}
	for _, dirName := range dirNames {
		files, err := os.ReadDir(filepath.Join(schemaDir, dirName))
		if err != nil {
			return retError("D_25", "Failed to lookup file", err)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if filepath.Ext(file.Name()) != ".sql" {
				continue
			}

			// obtain sql file
			schemaColumns, err := os.ReadFile(filepath.Join(schemaDir, dirName, file.Name()))
			if err != nil {
				return retError("D_27", "Failed to obtain schemaDir .sql file", err)
			}

			// obtain user file if it exists
			schemaUserColumns, err := os.ReadFile(filepath.Join(schemaUserDir, dirName, file.Name()))
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					writeWarning(fmt.Sprintf("DB: Detected that table %s has no defined user-specific columns",
						dirName+"/"+file.Name()))
				} else {
					return retError("D_28", "Failed to obtain schemaUserDir .sql file", err)
				}
			}

			// obtain table name
			tableName, err := grabTableName(string(schemaColumns), "sql")
			if err != nil {
				return retError("D_29", "Failed to obtain table name from .sql file", err)
			}

			// obtain include ID bool
			includeID, err := checkIncludeID(string(schemaColumns), "sql")
			if err != nil {
				return retError("D_30", "Failed to obtain include ID bool from .sql file", err)
			}

			// set correct metadata value
			var metadata string = userMetadata
			if includeID {
				metadata = userMetadataWithID
			}

			// create non-user table query
			nonUserQuery, err := generateTableQuery(tableName, string(schemaColumns), nonUserMetadata)
			if err != nil {
				// @TODO: impossible to reach
				return retError("D_31", "Failed to generate non-user table query", err)
			}

			// PokemonFormUser, PokemonSpeciesUser, etc.
			var userTableSuffix string = "User"
			var userTableName string = tableName + userTableSuffix

			// create user table query
			userColumns := string(schemaColumns)
			if len(schemaUserColumns) > 0 {
				sUCBody, sUCConstraints := grabConstraintsSQL(string(schemaUserColumns))
				userColumns = sUCBody + "\n" + userColumns + sUCConstraints
			}
			userQuery, err := generateTableQuery(userTableName, userColumns, metadata)
			if err != nil {
				// @TODO: impossible to reach
				return retError("D_32", "Failed to generate user table query", err)
			}

			_, err = tx.Exec(nonUserQuery)
			if err != nil {
				writeLog(nonUserQuery)
				message := fmt.Sprintf("Failed to create non-user table %s from %s; \n\toffending query outputted to error.txt\n\t",
					tableName,
					file.Name())
				return retError("D_33", message, err)
			}

			_, err = tx.Exec(userQuery)
			if err != nil {
				writeLog(fmt.Sprintf("%s\n\n%s", userColumns, userQuery))
				message := fmt.Sprintf("Failed to create user table %s from %s; \n\toffending query outputted to error.txt\n\t",
					userTableName,
					file.Name())
				return retError("D_34", message, err)
			}
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
	if err != nil {
		return nil, err
	}
	writeLog("DB: Created lookup tables")

	err = databaseSeedLookup(database)
	if err != nil {
		return nil, err
	}
	writeLog("DB: Seeded lookup tables")

	err = databaseCreateMain(database)
	if err != nil {
		return nil, err
	}
	writeLog("DB: Created main tables")

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

	resetMessageFile("logDB", "logDB.txt")
	writeFile(ret, "logDB.txt")
	return nil
}
