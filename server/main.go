package main

import (
	"os"
	"path"
	"photosync/src/database"
	"photosync/src/endpoint"
	"photosync/src/helper"
	"photosync/src/jwt"
	"photosync/src/metadata"
	"photosync/src/password"
	"runtime"

	"github.com/gin-gonic/gin"
)

func getDirectory() string {
	_, file, _, ok := runtime.Caller(1)
	if ok {
		return path.Dir(file)
	}

	return ""
}

func main() {
	db := database.NewPostgresDataBase("postgres", "postgres", "postgres", "localhost", 5432)
	passwordFacade := password.PasswordFacade{}
	timeHelper := helper.TimeHelper{}
	jwtManager := jwt.NewJwtManager(&timeHelper)
	rawMetadataExtractor := metadata.NewRawMetadataExtractor(getDirectory() + "/exiftool/exiftool")
	metadataExtractor := metadata.NewMetadataExtractor(&rawMetadataExtractor)
	hasher := helper.NewHasher()

	router := gin.Default()
	registerEndpoint := endpoint.NewRegisterEndpoint(db, passwordFacade)
	router.POST("/v1/register", registerEndpoint.Post)

	loginEndpoint := endpoint.NewLoginEndpoint(db, passwordFacade, &jwtManager, &timeHelper)
	router.POST("/v1/login", loginEndpoint.Post)

	uploadEndpoint := endpoint.NewUploadEndpoint(db, &metadataExtractor, &hasher, &jwtManager)
	router.POST("/v1/upload", uploadEndpoint.Post)

	fileDataEndpoint := endpoint.NewFileDataEndpoint(db, &jwtManager)
	router.GET("/v1/file_data", fileDataEndpoint.Get)

	fileEndpoint := endpoint.NewFileEndpoint(db, &jwtManager)
	router.GET("/v1/file", fileEndpoint.Get)
	router.DELETE("/v1/file", fileEndpoint.Delete)

	// TODO more FTies for delete !!!!

	if len(os.Args) == 2 && os.Args[1] == "--testing" {
		restartEndpoint := endpoint.NewRestartEndpoint()
		router.POST("/v1/restart", restartEndpoint.Post)
		router.HEAD("/v1/restart", restartEndpoint.Head)
	}

	router.Run()
}
