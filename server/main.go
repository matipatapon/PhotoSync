package main

import (
	"errors"
	"fmt"
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

func includeAllowedOriginInResponses(router *gin.Engine, envGetter helper.IEnvGetter) {
	allowedOrigin := envGetter.Get("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		panic(fmt.Errorf("'ALLOWED_ORIGIN' not specified"))
	}

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowedOrigin)
	})
}

func main() {
	envGetter := helper.NewEnvGetter()
	db, err := database.NewPostgresDataBase(&envGetter)
	if err != nil {
		panic(fmt.Errorf("error during database creation: '%s'", err.Error()))
	}

	err = db.InitDb()
	if err != nil {
		panic(fmt.Errorf("error during database initialization: '%s'", err.Error()))
	}

	passwordFacade := password.PasswordFacade{}
	timeHelper := helper.TimeHelper{}
	thumbnailCreator := helper.NewThumbnailCreator()
	jwtManager := jwt.NewJwtManager(&timeHelper)
	rawMetadataExtractor := metadata.NewRawMetadataExtractor(getDirectory() + "/exiftool/exiftool")
	metadataExtractor := metadata.NewMetadataExtractor(&rawMetadataExtractor)
	hasher := helper.NewHasher()

	router := gin.Default()

	includeAllowedOriginInResponses(router, &envGetter)

	registerEndpoint := endpoint.NewRegisterEndpoint(db, passwordFacade)
	router.POST("/v1/register", registerEndpoint.Post)

	loginEndpoint := endpoint.NewLoginEndpoint(db, passwordFacade, &jwtManager, &timeHelper)
	router.POST("/v1/login", loginEndpoint.Post)

	uploadEndpoint := endpoint.NewUploadEndpoint(db, &metadataExtractor, &hasher, &jwtManager, &thumbnailCreator)
	router.POST("/v1/upload", uploadEndpoint.Post)
	router.OPTIONS("/v1/upload", uploadEndpoint.Options)

	fileDataEndpoint := endpoint.NewFileDataEndpoint(db, &jwtManager)
	router.GET("/v1/file_data", fileDataEndpoint.Get)
	router.OPTIONS("/v1/file_data", fileDataEndpoint.Options)

	fileEndpoint := endpoint.NewFileEndpoint(db, &jwtManager)
	router.GET("/v1/file", fileEndpoint.Get)
	router.DELETE("/v1/file", fileEndpoint.Delete)
	router.OPTIONS("/v1/file", fileEndpoint.Options)

	datesEndpoint := endpoint.NewDatesEndpoint(db, &jwtManager)
	router.GET("/v1/dates", datesEndpoint.Get)
	router.OPTIONS("/v1/dates", datesEndpoint.Options)

	testing := envGetter.Get("TESTING")
	if testing == "true" {
		restartEndpoint := endpoint.NewRestartEndpoint(db)
		router.POST("/v1/restart", restartEndpoint.Post)
	} else if testing != "false" {
		panic(errors.New("'TESTING' has invalid value"))
	}

	isTlsEnabled := envGetter.Get("TLS_ENABLED")
	if isTlsEnabled == "true" {
		certPath := envGetter.Get("CERT_PATH")
		if certPath == "" {
			panic(errors.New("TLS enabled but 'CERT_PATH' not specified"))
		}
		certKeyPath := envGetter.Get("CERT_PRIVATE_KEY_PATH")
		if certKeyPath == "" {
			panic(errors.New("TLS enabled but 'CERT_PRIVATE_KEY_PATH' not specified"))
		}
		router.RunTLS(":8080", certPath, certKeyPath)
	} else if isTlsEnabled == "false" {
		router.Run()
	} else {
		panic(errors.New("TLS_ENABLED' has invalid value"))
	}
}
