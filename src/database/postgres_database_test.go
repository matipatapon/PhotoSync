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

func TestQueryShouldReturnErrorWhenCannotConnectToDb(t *testing.T) {
	sut := NewPostgresDataBase(getDb(), getUser(), "wrong_password", getIp(), getPort())

	result, err := sut.Query("SELECT version()")

	if len(result) != 0 || err == nil || !strings.HasPrefix(err.Error(), "failed to connect to ") {
		t.Error()
	}
}

func TestQueryShouldReturnResultOfQuery(t *testing.T) {
	db := createSut()

	result, err := db.Query("SELECT version()")

	if err != nil || !strings.HasPrefix(result[0][0].(string), "PostgreSQL") {
		t.Error()
	}
}

func TestQueryShouldReturnErrorWhenWrongQuery(t *testing.T) {
	db := createSut()

	result, err := db.Query("SELECT version)")

	if err == nil || len(result) != 0 {
		t.Error()
	}
}

func TestQueryShouldReturnErrorWhenTableHasNoRows(t *testing.T) {
	db := createSut()

	result, err := db.Query("SELECT * FROM postgres_database_test_empty_table")

	if err != nil || len(result) != 0 {
		t.Error(len(result))
	}
}

func TestQueryShouldReturnDataFromRow(t *testing.T) {
	db := createSut()

	result, err := db.Query("SELECT * FROM postgres_database_test_table_with_one_item")

	if err != nil || len(result) == 0 {
		t.Error()
	}
	if result[0][0] != int32(1) {
		t.Error()
	}
	if result[0][1] != "Mort" {
		t.Error()
	}
}

func TestQueryShouldReturnAllRows(t *testing.T) {
	db := createSut()

	result, err := db.Query("SELECT * FROM postgres_database_test_table_with_two_items ORDER BY id DESC")

	if err != nil || len(result) != 2 {
		t.Error()
	}
	if result[0][0] != int32(2) && result[0][1] != "Luna" {
		t.Error()
	}
	if result[1][0] != int32(1) && result[1][1] != "Mort" {
		t.Error()
	}
}

func TestExecuteShouldUpdateRecord(t *testing.T) {
	db := createSut()
	changedName := "changed_name"

	err := db.Execute("UPDATE postgres_database_test_table_to_update SET name = $1", changedName)

	if err != nil {
		t.Error()
	}
	result, err := db.Query("SELECT name FROM postgres_database_test_table_to_update")
	if len(result) != 1 || err != nil || result[0][0] != changedName {
		t.Error()
	}
}

func TestExecuteShouldDeleteRecord(t *testing.T) {
	db := createSut()

	err := db.Execute("DELETE FROM postgres_database_test_table_to_delete")

	if err != nil {
		t.Error()
	}
	result, err := db.Query("SELECT name FROM postgres_database_test_table_to_delete")
	if len(result) != 0 || err != nil {
		t.Error()
	}
}

func TestExecuteShouldInsertRecordsToDatabase(t *testing.T) {
	db := createSut()

	name1 := "Andrzej"
	id1 := 1
	name2 := "Grzegorz"
	id2 := 2

	err := db.Execute("INSERT INTO postgres_database_test_insertion_table VALUES($1, $2)", id1, name1)
	if err != nil {
		t.Error()
	}

	err = db.Execute("INSERT INTO postgres_database_test_insertion_table VALUES($1, $2)", id2, name2)
	if err != nil {
		t.Error()
	}

	result, err := db.Query("SELECT * FROM postgres_database_test_insertion_table ORDER BY name DESC")
	if len(result) != 2 || err != nil {
		t.Error()
	}

	if result[0][0] != id2 && result[0][1] != name2 {
		t.Error()
	}
	if result[1][0] != id1 && result[1][1] != name1 {
		t.Error()
	}

}

func TestExecuteShouldReturnErrorWhenInsertionFailed(t *testing.T) {
	db := createSut()

	err := db.Execute("INSERT INTO postgres_database_test_insertion_table VALUES($1, $2, $3)", 1, "name", "additional param")

	if err == nil {
		t.Error()
	}
}

func TestExecuteShouldReturnErrorWhenCannotConnectToDb(t *testing.T) {
	sut := NewPostgresDataBase(getDb(), getUser(), "wrong_password", getIp(), getPort())

	err := sut.Execute("SELECT version()")

	if err == nil || !strings.HasPrefix(err.Error(), "failed to connect to ") {
		t.Error()
	}
}
