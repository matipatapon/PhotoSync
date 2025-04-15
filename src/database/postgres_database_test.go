package database

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

func getDb() string {
	return os.Getenv("PGDB")
}

func getUser() string {
	return os.Getenv("PGUSER")
}

func getPassword() string {
	return os.Getenv("PGPASSWORD")
}

func getIp() string {
	return os.Getenv("PGIP")
}

func getPort() int {
	port, _ := strconv.Atoi(os.Getenv("PGPORT"))
	return port
}

func createSut() *PostgresDataBase {
	return NewPostgresDataBase(
		getDb(),
		getUser(),
		getPassword(),
		getIp(),
		getPort(),
	)
}

func TestQueryRowShouldReturnErrorWhenCannotConnectToDb(t *testing.T) {
	sut := NewPostgresDataBase(getDb(), getUser(), "wrong_password", getIp(), getPort())

	result, err := sut.QueryRow("SELECT version()")

	if result != nil || err == nil || !strings.HasPrefix(err.Error(), "failed to connect to ") {
		t.Error()
	}
}

func TestQueryRowShouldReturnResultOfQuery(t *testing.T) {
	db := createSut()

	result, err := db.QueryRow("SELECT version()")

	if err != nil || result == nil || !strings.HasPrefix(result[0].(string), "PostgreSQL") {
		t.Error()
	}
}

func TestQueryRowShouldReturnErrorWhenWrongQuery(t *testing.T) {
	db := createSut()

	result, err := db.QueryRow("SELECT version)")

	if err == nil || result != nil {
		t.Error()
	}
}

func TestQueryRowShouldReturnErrorWhenTableHasNoRows(t *testing.T) {
	db := createSut()

	result, err := db.QueryRow("SELECT * FROM postgres_database_test_empty_table")

	if err == nil || result != nil {
		t.Error()
	}
}

func TestQueryRowShouldReturnDataFromRow(t *testing.T) {
	db := createSut()

	result, err := db.QueryRow("SELECT * FROM postgres_database_test_table_with_one_item")

	if err != nil || result == nil {
		t.Error()
	}
	if result[0] != int32(1) {
		t.Error()
	}
	if result[1] != "Mort" {
		t.Error()
	}
}

func TestQueryRowShouldReturnFirstRow(t *testing.T) {
	db := createSut()

	result, err := db.QueryRow("SELECT * FROM postgres_database_test_table_with_two_items ORDER BY id DESC")

	if err != nil || result == nil {
		t.Error()
	}
	if result[0] != int32(2) {
		t.Error()
	}
	if result[1] != "Luna" {
		t.Error()
	}
}

// TODO write test for query with args
