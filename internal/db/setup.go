package db

/*
Database seeding and table construction code
*/

import (
	// "context"

	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/w7o/poke7db/internal/logging"

	_ "modernc.org/sqlite"
)

func seedTable(tx *sql.Tx, csvPath string, sqlPath string, tableIndex int, timestamp string) error {
	// Condition to exclude metadata
	includeMetadataCondition := tableIndex != 0

	csvContents, err := os.ReadFile(csvPath)
	if err != nil {
		return logging.RetError("D_11", "Failed to obtain .csv file", err)
	}

	sqlContents, err := os.ReadFile(sqlPath)
	if err != nil {
		message := fmt.Sprintf("csv -> sql file \"%s\" doesn't appear to exist", sqlPath)
		return logging.RetError("D_21", message, err)
	}

	tableName, err := grabTableName(string(csvContents), "csv")
	if err != nil {
		return logging.RetError("D_12", "Failed to grab table name csv", err)
	}

	includeID, err := checkIncludeID(string(sqlContents), "sql")
	if err != nil {
		return logging.RetError("D_20", "Failed to grab ID inclusion bool sql", err)
	}

	csvFile, _ := os.Open(csvPath)

	// Create a new reader for the CSV file
	reader := csv.NewReader(csvFile)
	reader.Comment = '#'

	// CSV Header
	columns, err := reader.Read()
	if err != nil {
		return logging.RetError("D_13", "Failed to read CSV header", err)
	}

	// Include createdAt timestamp if table contains metadata
	// NOTE: If editing to include more seeded metadata in the future, check
	// all includeMetadataCondition usages
	// %TODO make below variable a variable taht changes if includeID
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
	logging.WriteLog(fmt.Sprintf("DB: Query -  %s", query))
	statement, err := tx.Prepare(query)
	if err != nil {
		return logging.RetError("D_15", "Placeholder query preparation failed", err)
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
			return logging.RetError("D_14", "Failed to read CSV row", err)
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
			return logging.RetError("D_16", message, err)
		}

		args := convStringsToAnys(row)

		if includeMetadataCondition {
			// ..., createdAt, enabled
			args = append(args, timestamp, 1)
			if includeID {
				args = append(args, args[0])
			}
		}

		logging.WriteLog(fmt.Sprintf("DB: Arguments - %v", args))

		r, err := statement.Exec(args...)
		if err != nil {
			return logging.RetError("D_17", "Executing statement failed", err)
		}

		rowsAffected, err := r.RowsAffected()
		if err != nil {
			return logging.RetError("D_18", "Failed to get rows affected", err)
		}

		logging.WriteLog(fmt.Sprintf("DB: %d row[s] affected", rowsAffected))
	}
	return nil
}

/*
Index 00 is hardcoded to not have metadata applied
*/
func databaseCreateLookup(database *sql.DB) error {
	// Lookup tables have the same metadata as User Main tables
	logging.WriteLogAndConsole("DB: Creating lookup tables")
	err := databasePing(database)
	if err != nil {
		return err
	}

	tx, err := database.Begin()
	if err != nil {
		return logging.RetError("D_00", "", err)
	}
	defer tx.Rollback()

	// Create the lookup tables
	var dir string = "./database/schema_columns/00_lookup"
	// Golang ReadDir already sorts in alphabetical order
	files, err := os.ReadDir(dir)
	if err != nil {
		return logging.RetError("D_01", "", err)
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
			return logging.RetError("D_02", "Failed to extract filename index number", err)
		}
		if s == -1 {
			continue
		}

		// obtain the sql file
		schemaColumns, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return logging.RetError("D_03", "Failed to obtain .sql file", err)
		}

		// obtain table name from sql file
		tableName, err := grabTableName(string(schemaColumns), "sql")
		if err != nil {
			return logging.RetError("D_04", "Failed to obtain table name from .sql file", err)
		}

		// obtain include ID bool from sql file
		includeID, err := checkIncludeID(string(schemaColumns), "sql")
		if err != nil {
			return logging.RetError("D_22", "Failed to obtain include ID bool from .sql file", err)
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
			return logging.RetError("D_07", "Error when splitting body and constraint", err)
		}

		_, err = tx.Exec(query)
		if err != nil {
			logging.WriteLog(query)
			message := fmt.Sprintf("Failed to create table %s from %s; \n\toffending query outputted to error.txt\n\t",
				tableName, file.Name())
			return logging.RetError("D_05", message, err)
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
			return logging.RetError("D_26", "Failed to extract file index number", err)
		}
		if s == -1 {
			continue
		}

		var seedDir string = "./database/seed/00_lookup"
		csvPath := filepath.Join(seedDir, file.Name())
		csv, err := os.ReadFile(csvPath)
		if err != nil {
			return logging.RetError("D_35", "Failed to read CSV file", err)
		}

		tableName, err := grabTableName(string(csv), "csv")
		if err != nil {
			return logging.RetError("D_36", "Failed to read CSV table name", err)
		}

		clearQuery := fmt.Sprintf("DELETE FROM %s", tableName)
		_, err = tx.Exec(clearQuery)
		if err != nil {
			return logging.RetError("D_37", fmt.Sprintf("Failed to clear table %s",
				tableName), err)
		}
	}
	return nil
}

func databaseSeedLookup(database *sql.DB) error {
	logging.WriteLogAndConsole("DB: Seeding lookup tables")
	err := databasePing(database)
	if err != nil {
		return err
	}

	var seedDir string = "./database/seed/00_lookup"
	var schemaDir string = "./database/schema_columns/00_lookup"

	files, err := os.ReadDir(seedDir)
	if err != nil {
		return logging.RetError("D_09", "Reading seed lookup directory failed", err)
	}

	// set up database
	tx, err := database.Begin()
	if err != nil {
		return logging.RetError("D_08", "Database transaction init failed", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`PRAGMA defer_foreign_keys = ON`)
	if err != nil {
		return logging.RetError("D_38", "Failed to enable the deferring of foreign keys", err)
	}

	// @USER make this a setting/optional at some point
	err = clearLookupEntries(tx, files)
	if err != nil {
		return err
	}

	timestamp := time.Now().UTC().Format(time.RFC3339)

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
			return logging.RetError("D_10", "Failed to convert an index to int", err)
		}
		if s == -1 {
			continue
		}

		// obtain the csv file
		csvPath := filepath.Join(seedDir, file.Name())

		sqlPath := filepath.Join(schemaDir, sqlFileName)

		err = seedTable(tx, csvPath, sqlPath, s, timestamp)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func databaseCreateMain(database *sql.DB) error {
	logging.WriteLogAndConsole("DB: Creating main tables")
	err := databasePing(database)
	if err != nil {
		return err
	}

	tx, err := database.Begin()
	if err != nil {
		return logging.RetError("D_23", "Transaction start fail", err)
	}
	defer tx.Rollback()

	var schemaDir string = "./database/schema_columns"
	var schemaUserDir string = "./database/schema_columns_user"

	files, err := os.ReadDir(schemaDir)
	if err != nil {
		return logging.RetError("D_24", "Read schema col dir failed", err)
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
			return logging.RetError("D_25", "Failed to lookup file", err)
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
				return logging.RetError("D_27", "Failed to obtain schemaDir .sql file", err)
			}

			// obtain user file if it exists
			schemaUserColumns, err := os.ReadFile(filepath.Join(schemaUserDir, dirName, file.Name()))
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					logging.WriteWarning(fmt.Sprintf("DB: Detected that table %s has no defined user-specific columns",
						dirName+"/"+file.Name()))
				} else {
					return logging.RetError("D_28", "Failed to obtain schemaUserDir .sql file", err)
				}
			}

			// obtain table name
			tableName, err := grabTableName(string(schemaColumns), "sql")
			if err != nil {
				return logging.RetError("D_29", "Failed to obtain table name from .sql file", err)
			}

			// obtain include ID bool
			includeID, err := checkIncludeID(string(schemaColumns), "sql")
			if err != nil {
				return logging.RetError("D_30", "Failed to obtain include ID bool from .sql file", err)
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
				return logging.RetError("D_31", "Failed to generate non-user table query", err)
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
				return logging.RetError("D_32", "Failed to generate user table query", err)
			}

			_, err = tx.Exec(nonUserQuery)
			if err != nil {
				logging.WriteLog(nonUserQuery)
				message := fmt.Sprintf("Failed to create non-user table %s from %s; \n\toffending query outputted to error.txt\n\t",
					tableName,
					file.Name())
				return logging.RetError("D_33", message, err)
			}

			_, err = tx.Exec(userQuery)
			if err != nil {
				logging.WriteLog(fmt.Sprintf("%s\n\n%s", userColumns, userQuery))
				message := fmt.Sprintf("Failed to create user table %s from %s; \n\toffending query outputted to error.txt\n\t",
					userTableName,
					file.Name())
				return logging.RetError("D_34", message, err)
			}
		}
	}
	return tx.Commit()
}
