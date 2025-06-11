package main

import (
	"photosync/src/database"
	"photosync/src/endpoint"
	"photosync/src/helper"
	"photosync/src/jwt"
	"photosync/src/password"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.NewPostgresDataBase("postgres", "postgres", "postgres", "localhost", 5432)
	passwordFacade := password.PasswordFacade{}
	jwtManager := jwt.NewJwtManager()
	timeHelper := helper.TimeHelper{}

	router := gin.Default()
	registerEndpoint := endpoint.NewRegisterEndpoint(db, passwordFacade)
	router.POST("/register", registerEndpoint.Post)

	loginEndpoint := endpoint.NewLoginEndpoint(db, passwordFacade, &jwtManager, &timeHelper)
	router.GET("/login", loginEndpoint.Post)
	router.Run()
}
