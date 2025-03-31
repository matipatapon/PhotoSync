package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
)

type IDataBaseWrapper interface {
	QueryRow() ([]any, error)
}

type PostgresDataBaseWrapper struct {
	db       string
	user     string
	password string
	address  string
	port     int
}

func NewPostgresDataBaseWrapper(
	db string,
	user string,
	password string,
	address string,
	port int,
) *PostgresDataBaseWrapper {
	return &PostgresDataBaseWrapper{db, user, password, address, port}
}

func create_connection_url(
	db string,
	user string,
	password string,
	address string,
	port int,
) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, address, port, db)
}

func (dbw PostgresDataBaseWrapper) QueryRow(sql string) ([]any, error) {
	log.Printf("Connecting to %s", create_connection_url(dbw.db, dbw.user, "####", dbw.address, dbw.port))
	conn, err := pgx.Connect(context.Background(), create_connection_url(dbw.db, dbw.user, dbw.password, dbw.address, dbw.port))
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	log.Printf("Executing query %s", sql)
	rows, err := conn.Query(context.Background(), sql)

	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	x := rows.Next()

	log.Printf("%t", x)

	return rows.Values()
}
