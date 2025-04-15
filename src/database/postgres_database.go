package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
)

type PostgresDataBase struct {
	db       string
	user     string
	password string
	address  string
	port     int
}

func NewPostgresDataBase(
	db string,
	user string,
	password string,
	address string,
	port int,
) *PostgresDataBase {
	return &PostgresDataBase{db, user, password, address, port}
}

func (dbw PostgresDataBase) QueryRow(sql string, args ...any) ([]any, error) {
	log.Printf("Connecting to %s", createConnectionUrl(dbw.db, dbw.user, "####", dbw.address, dbw.port))
	conn, err := pgx.Connect(context.Background(), createConnectionUrl(dbw.db, dbw.user, dbw.password, dbw.address, dbw.port))
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	log.Printf("Executing query %s", sql)
	rows, err := conn.Query(context.Background(), sql, args...)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	defer rows.Close()
	rowReturned := rows.Next()
	log.Print("Fetching row")
	if !rowReturned {
		log.Print("No row returned!")
		return nil, errors.New("no rows returned")
	}

	return rows.Values()
}

func createConnectionUrl(
	db string,
	user string,
	password string,
	address string,
	port int,
) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, address, port, db)
}
