package data

// package to mock our db
// go get github.com/DATA-DOG/go-sqlmock

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	db2 "github.com/upper/db/v4"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	// set databasePool
	fakeDB, _, _ := sqlmock.New()
	defer fakeDB.Close()

	// run the func
	md := New(fakeDB)

	// set env var
	_ = os.Setenv("DATABASE_TYPE", "postgres")
	// assert Model struct, %T gives us the type
	if fmt.Sprintf("%T", md) != "data.Models" {
		t.Error("Wrong type", fmt.Sprintf("%T", md))
	}

	// do the same for mysql
	_ = os.Setenv("DATABASE_TYPE", "mysql")
	if fmt.Sprintf("%T", md) != "data.Models" {
		t.Error("Wrong type", fmt.Sprintf("%T", md))
	}
}

func TestGetInsertID(t *testing.T) {
	// test the id back from postgres which is int64
	// Arrange
	var id db2.ID
	id = int64(1)

	// Act
	res := GetInsertID(id)

	// Assert
	if fmt.Sprintf("%T", res) != "int" {
		t.Error("Wrong type", fmt.Sprintf("%T", res))
	}

	// test the id back from mariaDB/mysql which is int
	// test2
	id = 1
	// Act
	res = GetInsertID(id)

	// Assert
	if fmt.Sprintf("%T", res) != "int" {
		t.Error("Wrong type", fmt.Sprintf("%T", res))
	}

}
