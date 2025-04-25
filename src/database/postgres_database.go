package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

var logger *log.Logger = log.New(os.Stdout, "[PostgresDataBase]: ", log.LstdFlags)

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

// Query method overrides IDataBase.Query.
// Error will be returned when:
//   - Connection to db cannot be established
//   - PostgreSQL database will return error when processing query
func (dbw PostgresDataBase) Query(sql string, args ...any) ([][]any, error) {
	logger.Printf("Connecting to %s", createConnectionUrl(dbw.db, dbw.user, "####", dbw.address, dbw.port))
	conn, err := pgx.Connect(context.Background(), createConnectionUrl(dbw.db, dbw.user, dbw.password, dbw.address, dbw.port))
	if err != nil {
		logger.Print(err.Error())
		return nil, err
	}

	logger.Printf("Executing query %s", sql)
	rows, err := conn.Query(context.Background(), sql, args...)
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
