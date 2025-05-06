package main

import (
	"photosync/src/database"
	"photosync/src/endpoint"

	"github.com/gin-gonic/gin"
)

var ZIEMNIAK = false

type RegisterData struct {
	username string
	password string
}

func register(c *gin.Context) {
	registerData := RegisterData{}
	c.BindJSON(&registerData)
	c.Status(200)

	print(registerData.password, registerData.username)
}

func main() {
	db := database.NewPostgresDataBase("postgres", "postgres", "postgres", "localhost", 5432)
	db.Query("SELECT version()")
	router := gin.Default()
	registerEndpoint := endpoint.RegisterEndpoint{}
	router.POST("/register", registerEndpoint.Post)
	router.Run()
}
