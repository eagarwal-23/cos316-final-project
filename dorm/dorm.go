package dorm

import (
	"database/sql"
	"log"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

// DB handle
type DB struct {
	inner *sql.DB
}

type Student struct {
	Name      string
	Age       int
	ClassYear int
	Major     string
}

// NewDB returns a new DB using the provided `conn`,
// an sql database connection.
// This function is provided for you. You DO NOT need to modify it.
func NewDB(conn *sql.DB) DB {
	return DB{inner: conn}
}

// Close closes db's database connection.
// This function is provided for you. You DO NOT need to modify it.
func (db *DB) Close() error {
	return db.inner.Close()
}

// ColumnNames analyzes a struct, v, and returns a list of strings,
// one for each of the public fields of v.
// The i'th string returned should be equal to the name of the i'th
// public field of v, converted to underscore_case.
// Refer to the specification of underscore_case, below.

// Example usage:
// type MyStruct struct {
//    ID int64
//    UserName string
// }
// ColumnNames(&MyStruct{})    ==>   []string{"id", "user_name"}
func ColumnNames(v interface{}) []string {
	columnNames := make([]string, 0)

	val := reflect.ValueOf(v)

	// while val is a pointer, resolve its value
	for val.Kind() == reflect.Ptr {

		val = reflect.Indirect(val)
	}

	// should double check we now have a struct (could still be anything)
	if val.Kind() != reflect.Struct {
		panic("unexpected type")
	}

	// now we grab our values
	structType := val.Type()
	numFields := structType.NumField()

	for i := 0; i < numFields; i++ {
		varName := structType.Field(i).Name
		if unicode.IsUpper(rune(varName[0])) {
			columnNames = append(columnNames, camelToUnderscore(varName))
		}

	}
	return columnNames
}

// TableName analyzes a struct, v, and returns a single string, equal
// to the name of that struct's type, converted to underscore_case.
// Refer to the specification of underscore_case, below.

// Example usage:
// type MyStruct struct {
//    ...
// }
// TableName(&MyStruct{})    ==>  "my_struct"
func TableName(result interface{}) string {
	resultString := reflect.TypeOf(result).String()
	resultArray := strings.Split(resultString, ".")
	resultString = strings.TrimLeft(resultString, resultArray[0])
	resultString = strings.TrimLeft(resultString, ".")

	return camelToUnderscore(resultString)
}

// Find queries a database for all rows in a given table,
// and stores all matching rows in the slice provided as an argument.
// Find should panic if the table doesn't exist.

// The argument `result` will be a pointer to an empty slice of models.
// To be explicit, it will have type: *[]MyStruct,
// where MyStruct is any arbitrary struct subject to the restrictions
// discussed later in this document.
// You may assume the slice referenced by `result` is empty.

// Example usage to find all UserComment entries in the database:
//    type UserComment struct = { ... }
//    result := []UserComment{}
//    db.Find(&result)

// https://stackoverflow.com/questions/31924199/get-fields-of-empty-struct-slice-in-go
// https://stackoverflow.com/questions/69868784/append-to-golang-slice-passed-as-empty-interface
func (db *DB) Find(result interface{}) {
	// Get name of table
	tableName := TableName(result)

	resultStruct := reflect.ValueOf(result).Elem().Type().Elem()

	fields := make([]reflect.StructField, resultStruct.NumField())
	for i, _ := range fields {
		fields[i] = resultStruct.Field(i)
	}

	structType := reflect.StructOf(fields)

	stmt, err := db.inner.Prepare("SELECT * FROM " + tableName)
	if err != nil {
		panic(err)
	}

	rows, err := stmt.Query()

	for rows.Next() {
		varStruct := reflect.New(structType)

		addresses := make([]interface{}, 0)
		for i := 0; i < len(fields); i++ {
			addresses = append(addresses, reflect.Indirect(varStruct).Field(i).Addr().Interface())
		}

		if err := rows.Scan(addresses...); err != nil {
			log.Fatal(err)
		}

		tempResult := reflect.ValueOf(result).Elem()
		a := reflect.Append(tempResult, reflect.Indirect(varStruct))
		reflect.ValueOf(result).Elem().Set(a)
	}
}

// First queries a database for the first row in a table,
// and stores the matching row in the struct provided as an argument.
// If no such entry exists, First returns false; else it returns true.
// First should panic if the table doesn't exist.

// The argument `result` will be a pointer to a model.
// To be explicit, it will have type: *MyStruct,
// where MyStruct is any arbitrary struct subject to the restrictions
// discussed later in this document.

// Example usage to find the first UserComment entry in the database:
//    type UserComment struct = { ... }
//    result := &UserComment{}
//    ok := db.First(result)
// with the argument), otherwise return true.
func (db *DB) First(result interface{}) bool {
	tableName := TableName(result)
	columnNames := ColumnNames(result)
	queryString := strings.Join(columnNames, ", ")

	stmt, err := db.inner.Prepare("SELECT " + queryString + " FROM " + tableName + " LIMIT 1")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	addresses := scannerStruct(result)
	err = stmt.QueryRow().Scan(addresses...)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case of no rows returned.
			return false
		}
		panic(err)
	}
	return true
}

// Create adds the specified model to the appropriate database table.
// The table for the model *must* already exist, and Create() should
// panic if it does not.

// The argument `model` will be a pointer to a model.
// To be explicit, it will have type *MyStruct,
// where MyStruct is any arbitrary struct subject to the restrictions
// discussed later in this document.

// Optionally, at most one of the fields of the provided `model`
// might be annotated with the tag `dorm:"primary_key"`. If such a
// field exists, Create() should ignore the provided value of that
// field, overwriting it with the auto-incrementing row ID.
// This ID is given by the value of last_inserted_rowid(),
// returned from the underlying sql database.

// Example usage to add a new user row to the database:
//    type User struct = {
//        Name string
//        ID   int64 `dorm:"primary_key"`
//    }
//    user := &User{Name: "NAME", ID: 20}
//    db.Create(user)
//    user.ID should now be updated to last_inserted_rowid()
func (db *DB) Create(model interface{}) {
	var primaryPresent bool                // true if table has a primary key
	var primaryFieldNo int				   // field no. of primary key 
	insertColumns := make([]string, 0)     // store columns that are updated 
	insertValues := make([]interface{}, 0) // store values to insert into table
	tableName := TableName(model)

	// Extract values from `model`: indirect because this is a pointer to a model
	val := reflect.ValueOf(model)
	val = reflect.Indirect(val)
	valType := reflect.TypeOf(model).Elem()
	numFields := val.NumField()

	// Populate insertColumns and insertValues
	for i := 0; i < numFields; i++ {
		field := val.Field(i)
		tag := valType.Field(i).Tag

		// Ignore field if it is primary key
		if tag != "dorm:\"primary_key\"" {
			column := camelToUnderscore(val.Type().Field(i).Name)
			insertColumns = append(insertColumns, column)
			insertValues = append(insertValues, field.Interface())
		} else if tag == "dorm:\"primary_key\"" {
			primaryPresent = true
			primaryFieldNo = i
		}

	}

	// DATABASE OPERATIONS
	queryString := strings.Join(insertColumns, ", ")

	// if primary key is present in model we have one less field to update
	if primaryPresent {
		numFields = numFields - 1 
	}

	stmt, err := db.inner.Prepare("INSERT INTO " + tableName + "(" + queryString + ")" + "VALUES(" + createPlaceholders(numFields) + ")")
	if err != nil {
		panic(err) // panic if table doesn't exist
	}

	// Insert to table.
	res, err := stmt.Exec(insertValues...)
	if err != nil {
		panic(err)
	}

	// update ID to last_inserted_rowid()
	lastId, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	reflect.ValueOf(model).Elem().Field(primaryFieldNo).SetInt(lastId)
}

// Helper function to convert variable names from camel_Case to underscore_case
func camelToUnderscore(camelString string) string {
	re := regexp.MustCompile(`(.)([A-Z][a-z])`)
	re2 := regexp.MustCompile(`([a-z0-9])([A-Z])`)

	underscoreString := re.ReplaceAllString(camelString, "${1}_${2}")
	underscoreString = re2.ReplaceAllString(underscoreString, "${1}_${2}")

	return strings.ToLower(underscoreString)
}

// Helper function to create struct to pass into db.Scan
func scannerStruct(u interface{}) []interface{} {
	valT := reflect.ValueOf(u)
	for valT.Kind() == reflect.Ptr {
		valT = reflect.Indirect(valT)
	}

	numFields := valT.NumField()
	addresses := make([]interface{}, 0)

	for i := 0; i < numFields; i++ {
		addresses = append(addresses, valT.Field(i).Addr().Interface())
	}

	return addresses
}

// Helper function to create db placeholders
func createPlaceholders(numFields int) string {
	placeHolder := ""
	for i := 0; i < numFields; i += 1 {
		placeHolder += "?,"
	}
	placeHolder = strings.TrimRight(placeHolder, ",")
	return placeHolder
}
