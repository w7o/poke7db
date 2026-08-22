package db

import (
	"database/sql"
	"fmt"
	"reflect"

	// "strconv"
	"strings"
	"time"

	"github.com/w7o/poke7db/internal/logging"
)

func nullableString(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

// %MDO: structs instead of managing name and value separately..?

// Inserts or updates the given data into a database --
// dbStruct MUST be a list of THE SAME struct
func upsertStruct(tx *sql.Tx, tableName string, dbStruct any,
	metadata metadataTemplate) error {

	// Get a reflect.Value representing a struct list's value
	structReflect := reflect.ValueOf(dbStruct)

	if structReflect.Kind() != reflect.Slice {
		//%TODO
		return logging.RetError("F_01",
			fmt.Sprintf("Passed in type %T is not a list", dbStruct),
			nil)
	}

	numRows := structReflect.Len()
	if numRows == 0 {
		return logging.RetError("F_02", "Passed in list doesn't contain anything", nil)
	}

	// Check if the value of the first item in the list is a struct, and that
	// each value in the list is the same type
	structType := structReflect.Index(0).Type()
	if structType.Kind() != reflect.Struct {
		return logging.RetError("F_03", "Passed in list isn't made of structs", nil)
	}
	for i := range numRows {
		if structReflect.Index(i).Type() != structType {
			return logging.RetError("F_08",
				"All values in list are not the same struct, or the list isn't entirely made of structs",
				nil)
		}
	}

	// Returns the number of numFields in the struct
	numFields := structType.NumField()

	// if need DB tag implementation: sDC 719-06#4
	var fieldNames []string
	var excludedFieldNames []string
	var fieldValues []any
	var placeholders []string

	// grab primary key sql.Rows
	pkQuery := fmt.Sprintf(
		"SELECT name FROM pragma_table_info('%s') WHERE pk > 0 ORDER BY pk",
		tableName)
	pkRows, err := tx.Query(pkQuery)
	if err != nil {
		logging.WriteLog(pkQuery)
		return logging.RetError("F_06",
			"Primary key query failed; offending query written to log", err)
	}
	defer pkRows.Close()

	// extract primary keys from sql.Rows into excludeFieldNames
	primaryKeys, err := convArbitraryRowsToStrings(pkRows, 1)
	if err != nil {
		return logging.RetError("F_07",
			"Converting extracted rows.Sql into string failed", err)
	}
	if len(primaryKeys) == 0 {
		return logging.RetError("F_10", "Table has no primary key", nil)
	}
	excludedFieldNames = append(excludedFieldNames, primaryKeys...)

	// returns all table column names minus metadata
	for i := range numFields {
		fieldNames = append(fieldNames, structType.Field(i).Name)
	}

	// check imported metadata field validity:
	// if importedAt is true then it's a non-user table
	hasNonUser := metadata.importedAt != nil
	// if createdAt is true then it's a user table
	hasUser := metadata.createdAt != nil
	// ID can only exist if it's a user table
	idConditionFail := !hasUser && metadata.hasID
	if (hasNonUser == hasUser) || (idConditionFail) {
		return logging.RetError("F_05", "Invalid metadata field combination", nil)
	}

	var metadataValues []any

	// handling metadata
	fieldNames = append(fieldNames, "originID", "enabled")
	excludedFieldNames = append(excludedFieldNames, "enabled")
	metadataValues = append(metadataValues,
		metadata.originID, //originID
		1,                 //enabled
	)
	if hasNonUser {
		fieldNames = append(fieldNames, "importedAt", "checkedAt")
		excludedFieldNames = append(excludedFieldNames, "importedAt")
		metadataValues = append(metadataValues,
			nullableString(metadata.importedAt), //importedAt
			nullableString(metadata.checkedAt),  //checkedAt
		)
	} else { // hasUser
		fieldNames = append(fieldNames, "createdAt", "updatedAt")
		excludedFieldNames = append(excludedFieldNames, "createdAt")
		metadataValues = append(metadataValues,
			nullableString(metadata.createdAt), //createdAt
			nullableString(metadata.updatedAt), //updatedAt
		)
		if metadata.hasID {
			fieldNames = append(fieldNames, "ID")
			excludedFieldNames = append(excludedFieldNames, "ID")
			// metadata value insertion done in below loop
		}
	}

	logging.WriteLog(fmt.Sprintf("DB: Field names: %v", fieldNames))
	logging.WriteLog(fmt.Sprintf("DB: Excluded fields: %v", excludedFieldNames))

	// construct placeholders
	var placeholderIndex int = 1
	for i := 0; i < structReflect.Len(); i++ {
		// currStruct represents the current struct (i.e. table)
		currStruct := structReflect.Index(i)
		var rowPlaceholders []string

		// add placeholders for the current row
		for j := range numFields {
			fieldValues = append(fieldValues, currStruct.Field(j).Interface())
			rowPlaceholders = append(rowPlaceholders,
				fmt.Sprintf("$%d", placeholderIndex))
			placeholderIndex++
		}

		// add metadata values
		for _, value := range metadataValues {
			fieldValues = append(fieldValues, value)
			rowPlaceholders = append(rowPlaceholders,
				fmt.Sprintf("$%d", placeholderIndex))
			placeholderIndex++
		}

		// add ID w/ value of primary key if it's enabled
		if metadata.hasID {
			if len(primaryKeys) > 1 {
				return logging.RetError("F_09",
					"ID is enabled when there are multiple primary keys; cannot derive ID value",
					nil)
			}
			fieldValues = append(
				fieldValues,
				currStruct.Field(0).Interface(),
			)
			rowPlaceholders = append(rowPlaceholders,
				fmt.Sprintf("$%d", placeholderIndex))
			placeholderIndex++
		}

		placeholders = append(placeholders,
			"("+strings.Join(rowPlaceholders, ", ")+")",
		)
	}

	// build on conflict update section
	var updateFields []string
	for _, fieldName := range fieldNames {
		excluded := false

		for _, excludedField := range excludedFieldNames {
			if fieldName == excludedField {
				excluded = true
				break
			}
		}

		if !excluded {
			updateFields = append(updateFields,
				fmt.Sprintf("%s = EXCLUDED.%s", fieldName, fieldName))
		}
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s ON CONFLICT (%s) DO UPDATE SET %s",
		tableName,
		strings.Join(fieldNames, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(primaryKeys, ", "),
		strings.Join(updateFields, ", "),
	)

	_, err = tx.Exec(query, fieldValues...)
	logging.WriteLog(fmt.Sprintf("DB:\n\tquery:%s\n\tvalues: %v", query, fieldValues))
	if err != nil {
		return logging.RetError("F_04", "Execution of query failed; offending query saved in log.txt", err)
	}
	return nil
}

func upsertTables(tables []tableStruct, originID Origin) error {
	// %TODO fix call to this and also finish this
	tx, err := database.Begin()
	if err != nil {
		return logging.RetError("F_00", "Error beginning database transaction", err)
	}
	defer tx.Rollback()

	// %TODO
	// Set up metadata (currently only for non-user tables)
	timestamp := time.Now().UTC().Format(time.RFC3339)
	metadata := metadataTemplate{
		originID:   int(originID),
		importedAt: &timestamp,
		checkedAt:  &timestamp,
		enabled:    1,
		hasID:      false,
	}

	for _, table := range tables {
		logging.WriteLog("===\nBEFORE")
		err = d_queryAndLog(tx, d_testQueries[table.TableName]+d_dexNum)
		if err != nil {
			return err
		}
		err = upsertStruct(tx, table.TableName, table.TableVar, metadata)
		if err != nil {
			return err
		}
		logging.WriteLog("AFTER")
		err = d_queryAndLog(tx, d_testQueries[table.TableName]+d_dexNum)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
