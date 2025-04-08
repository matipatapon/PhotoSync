package database

import (
	"strings"
	"testing"
)

func Test_PostgresQueryRow_ShouldReturnErrorWhenCannotConnectToDb(t *testing.T) {
	db := NewPostgresDataBase("postgres", "postgres", "pop", "localhost", 5432)
	result, err := db.QueryRow("SELECT version()")
	if result != nil || err == nil || !strings.HasPrefix(err.Error(), "failed to connect to ") {
		t.Error()
	}
}

func Test_PostgresQueryRow_ShouldReturnResultOfQuery(t *testing.T) {
	db := NewPostgresDataBase("postgres", "postgres", "postgres", "localhost", 5432)
	result, err := db.QueryRow("SELECT version()")
	if err != nil || result == nil || !strings.HasPrefix(result[0].(string), "PostgreSQL") {
		t.Error()
	}
}

func Test_PostgresQueryRow_ShouldReturnErrorWhenWrongQuery(t *testing.T) {
	db := NewPostgresDataBase("postgres", "postgres", "postgres", "localhost", 5432)
	result, err := db.QueryRow("SELECT version)")
	if err == nil || result != nil {
		t.Error()
	}
}
