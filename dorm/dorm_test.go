package dorm

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func connectSQL() *sql.DB {
	conn, err := sql.Open("sqlite3", "file:test.db?mode=memory")
	if err != nil {
		panic(err)
	}
	return conn
}

func createUserTable(conn *sql.DB) {
	_, err := conn.Exec(`create table user ( 
		user_id INTEGER PRIMARY KEY AUTOINCREMENT, full_name text, year int, bool BOOL, float FLOAT, rune CHAR, byte int
	)`)

	if err != nil {
		panic(err)
	}
}

func insertUsers(conn *sql.DB, users []User) {
	for _, uc := range users {
		_, err := conn.Exec(`insert into user
		values
		(?, ?, ?, ?, ?, ?, ?)`,uc.UserID, uc.FullName, uc.Year, uc.Bool, uc.Float, uc.Rune, uc.Byte)

		if err != nil {
			panic(err)
		}
	}
}

// type User struct {
// 	UserID   int `dorm:"primary_key"`
// 	FullName string
// 	Year     int 
// }

type User struct {
	UserID   int `dorm:"primary_key"`
	FullName string
	Year     int
	Bool     bool
	Float    float32
	Rune     rune
	Byte     byte
}

type StudentInfo struct {
	Name      string
	Age       int
	ClassYear int
	Major     string
}

// var MockUsers = []User{
// 	User{FullName: "Test User1", Year: 1990},
// 	User{FullName: "Test User2", Year: 2000},
// 	User{FullName: "Test User3", Year: 2049},
// 	User{FullName: "Test User4", Year: 1346},
// }

var MockUsers = []User{
	// User{FullName: "Test User1", Year: 1990},
	// User{FullName: "Test User2", Year: 2000},
	// User{FullName: "Test User3", Year: 2049},
	User{UserID: 1, FullName: "Test User1", Year: 1990, Bool: true, Float: 1.2, Rune: 'B', Byte: 97},
}

// func TestFind(t *testing.T) {
// 	conn := connectSQL()
// 	createUserTable(conn)
// 	insertUsers(conn, MockUsers)

// 	db := NewDB(conn)
// 	defer db.Close()

// 	results := []User{}
// 	//p := &results
// 	fmt.Println("the result is prev", results)

// 	db.Find(&results)
// 	fmt.Println("the result is after", results)
// 	//fmt.Printf("the result in test is %v\n", results)

// 	if len(results) != 3 {
// 		t.Errorf("Expected 3 users but found %d", len(results))
// 	}
// }

func TestColumnNames(t *testing.T) {
	outputArray := ColumnNames(&StudentInfo{})
	checkArray := []string{"name", "age", "class_year", "major"}
	fmt.Printf("%v", outputArray)
	fmt.Printf("TestColumnNames %v\n", Equal(outputArray, checkArray))
}

func TestTableName(t *testing.T) {
	outputString := TableName(&StudentInfo{})
	checkString := "student_info"
	fmt.Printf("outputString %v \n", outputString)
	fmt.Printf("checkString %v \n", checkString)
	fmt.Printf("TestTableName %v\n", strings.Compare(checkString, outputString) == 0)
}

// func TestFirst(t *testing.T) {
// 	conn := connectSQL()
// 	createUserTable(conn)
// 	insertUsers(conn, MockUsers)

// 	db := NewDB(conn)
// 	defer db.Close()

// 	result := &User{}
// 	ok := db.First(result)

// 	if !ok {
// 		t.Errorf("Expected 1 users but found %v", ok)
// 	}
// }

func TestCreate(t *testing.T) {
	conn := connectSQL()
	createUserTable(conn)
	insertUsers(conn, MockUsers)

	db := NewDB(conn)
	defer db.Close()

	user1 := &User{FullName: "Test UserX", Year: 1940, Bool: false, Float: 3.4, Rune: 'X', Byte: 98}
	user2 := &User{FullName: "Test UserY", Year: 1942, Bool: true, Float: 3.9, Rune: 'Y', Byte: 99}
	user3 := &User{FullName: "Test UserZ", Year: 2023, Bool: true, Float: 3.9, Rune: 'Z', Byte: 100}
	db.Create(user1)
	db.Create(user2)
	db.Create(user3)

	//t.Errorf("TESTING CREATE")
}

func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
