package main

import (
	"os"
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
	timeHelper := helper.TimeHelper{}
	jwtManager := jwt.NewJwtManager(&timeHelper)

	router := gin.Default()
	registerEndpoint := endpoint.NewRegisterEndpoint(db, passwordFacade)
	router.POST("/v1/register", registerEndpoint.Post)

	loginEndpoint := endpoint.NewLoginEndpoint(db, passwordFacade, &jwtManager, &timeHelper)
	router.POST("/v1/login", loginEndpoint.Post)

	uploadEndpoint := endpoint.NewUploadEndpoint()
	router.POST("/v1/upload", uploadEndpoint.Post)

	if len(os.Args) == 2 && os.Args[1] == "--testing" {
		exitEndpoint := endpoint.NewExitEndpoint()
		router.POST("/v1/exit", exitEndpoint.Post)
	}

	router.Run()
}
