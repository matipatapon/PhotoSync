package main

import (
	"photosync/src/database"

	"github.com/gin-gonic/gin"
)

var ZIEMNIAK = false

func main() {
	db := database.NewPostgresDataBase("postgres", "postgres", "postgres", "localhost", 5432)
	db.Query("SELECT version()")
	router := gin.Default()
	router.Run()
}
