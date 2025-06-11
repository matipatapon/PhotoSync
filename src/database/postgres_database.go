package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
)

var logger *log.Logger = log.New(os.Stdout, "[PostgresDataBase]: ", log.LstdFlags)

var TIMEOUT time.Duration = time.Second * 30

// PostgresDataBase struct implements IDataBase interface.
// It handles connection with PostgreSQL database.
type PostgresDataBase struct {
	db       string
	user     string
	password string
	address  string
	port     int
}

// NewPostgresDataBase function creates PostgresDataBase.
func NewPostgresDataBase(
	db string,
	user string,
	password string,
	address string,
	port int,
) *PostgresDataBase {
	return &PostgresDataBase{db, user, password, address, port}
}

// Execute method overrides IDataBase.Query.
// Error will be returned when:
//   - Connection to db cannot be established
//   - PostgreSQL database will return error when processing query
func (dbw PostgresDataBase) Execute(sql string, args ...any) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	logger.Printf("Connecting to '%s'", createConnectionUrl(dbw.db, dbw.user, "####", dbw.address, dbw.port))
	conn, err := pgx.Connect(ctx, createConnectionUrl(dbw.db, dbw.user, dbw.password, dbw.address, dbw.port))
	if err != nil {
		logger.Print(err.Error())
		return err
	}
	defer conn.Close(ctx)

	logger.Printf("Executing modifying query '%s'", sql)
	_, err = conn.Exec(ctx, sql, args...)
	if err != nil {
		logger.Printf("Execution failed %s", err.Error())
	}
	return err
}

// Query method overrides IDataBase.Query.
// Error will be returned when:
//   - Connection to db cannot be established
//   - PostgreSQL database will return error when processing query
func (dbw PostgresDataBase) Query(sql string, args ...any) ([][]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	logger.Printf("Connecting to '%s'", createConnectionUrl(dbw.db, dbw.user, "####", dbw.address, dbw.port))
	conn, err := pgx.Connect(ctx, createConnectionUrl(dbw.db, dbw.user, dbw.password, dbw.address, dbw.port))
	if err != nil {
		logger.Print(err.Error())
		return nil, err
	}
	defer conn.Close(ctx)

	logger.Printf("Executing non-modifying query '%s'", sql)
	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		logger.Print(err.Error())
		return nil, err
	}
	defer rows.Close()

	result := [][]any{}
	logger.Print("Getting rows ...")
	count := 0
	for {
		if !rows.Next() {
			break
		}
		row, _ := rows.Values()
		result = append(result, row)
		count += 1
	}
	logger.Printf("Got %d rows", count)
	return result, nil
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
