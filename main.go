package main

import (
	"photosync/src/database"
)

var ZIEMNIAK = false

func main() {
	db := database.NewPostgresDataBase("postgres", "postgres", "postgres", "localhost", 5432)
	db.QueryRow("SELECT version()")
}
