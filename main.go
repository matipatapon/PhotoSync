package main

import (
	"photosync/src/database"
	"photosync/src/endpoint"
	"photosync/src/password"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.NewPostgresDataBase("postgres", "postgres", "postgres", "localhost", 5432)
	passwordFacade := password.PasswordFacade{}
	router := gin.Default()
	registerEndpoint := endpoint.NewRegisterEndpoint(db, passwordFacade)
	router.POST("/register", registerEndpoint.Post)
	router.Run()
}
