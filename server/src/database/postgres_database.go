package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"photosync/src/helper"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
)

var logger *log.Logger = log.New(os.Stdout, "[PostgresDataBase]: ", log.LstdFlags)

var TIMEOUT time.Duration = time.Second * 30

type PostgresDataBase struct {
	db       string
	user     string
	password string
	address  string
	port     int
}

func NewPostgresDataBase(
	getter helper.IEnvGetter) (*PostgresDataBase, error) {

	db := getter.Get("DATABASE")
	if db == "" {
		return nil, errors.New("env 'DATABASE' is missing")
	}

	user := getter.Get("DATABASE_USER")
	if user == "" {
		return nil, errors.New("env 'DATABASE_USER' is missing")
	}

	password := getter.Get("DATABASE_PASSWORD")
	if password == "" {
		return nil, errors.New("env 'DATABASE_PASSWORD' is missing")
	}

	address := getter.Get("DATABASE_IP")
	if address == "" {
		return nil, errors.New("env 'DATABASE_IP' is missing")
	}

	portRaw := getter.Get("DATABASE_PORT")
	if portRaw == "" {
		return nil, errors.New("env 'DATABASE_PORT' is missing")
	}

	port, err := strconv.ParseInt(portRaw, 10, 32)
	if err != nil {
		return nil, err
	}

	return &PostgresDataBase{db, user, password, address, int(port)}, nil
}

func (dbw PostgresDataBase) InitDb() error {
	err := dbw.Execute(`
		CREATE TABLE IF NOT EXISTS users(
			id bigserial,
			username text NOT NULL,
			password text NOT NULL,
			PRIMARY KEY (id),
			CONSTRAINT "username is unique" UNIQUE (username)
		);
		CREATE TABLE IF NOT EXISTS files(
			id bigserial,
			user_id bigint REFERENCES users(id) NOT NULL,
			creation_date timestamp NOT NULL,
			filename text NOT NULL,
			mime_type smallint NOT NULL,
			file bytea NOT NULL,
			hash text NOT NULL,
			size bigint NOT NULL,
			PRIMARY KEY (id),
			CONSTRAINT "file is unique" UNIQUE (user_id, hash, size)
		);
	`)
	if err != nil {
		logger.Printf("Failed to initialize db: '%s'", err.Error())
	}
	return err
}

func (dbw PostgresDataBase) DropDb() error {
	err := dbw.Execute(`
		DROP TABLE files;
		DROP TABLE users;
	`)
	if err != nil {
		logger.Printf("Failed to drop db: '%s'", err.Error())
	}
	return err
}

func (dbw PostgresDataBase) Execute(sql string, args ...any) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	conn, err := pgx.Connect(ctx, createConnectionUrl(dbw.db, dbw.user, dbw.password, dbw.address, dbw.port))
	if err != nil {
		logger.Print(err.Error())
		return err
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, sql, args...)
	if err != nil {
		logger.Printf("Execution failed %s", err.Error())
		return err
	}

	logger.Printf("Executed modifying query '%s'", sql)
	return err
}

func (dbw PostgresDataBase) Query(sql string, args ...any) ([][]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	conn, err := pgx.Connect(ctx, createConnectionUrl(dbw.db, dbw.user, dbw.password, dbw.address, dbw.port))
	if err != nil {
		logger.Print(err.Error())
		return nil, err
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		logger.Print(err.Error())
		return nil, err
	}
	defer rows.Close()

	result := [][]any{}
	count := 0
	for {
		if !rows.Next() {
			break
		}
		row, _ := rows.Values()
		result = append(result, row)
		count += 1
	}
	logger.Printf("Executed non-modifying query '%s', returned '%d' rows", sql, count)

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
