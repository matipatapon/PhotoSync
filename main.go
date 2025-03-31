package main

import (
	"photosync/src/database"
)

var ZIEMNIAK = false

func main() {
	db := database.NewPostgresDataBaseWrapper("postgres", "postgres", "postgres", "localhost", 5432)
	db.QueryRow("SELECT * from Test")
}
