package main

import (
	"errors"
	"fmt"
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

func runTLSIfEnabled(router *gin.Engine, envGetter helper.IEnvGetter) {
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
	} else if isTlsEnabled != "false" {
		panic(errors.New("TLS_ENABLED' has invalid value"))
	}
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
		panic(fmt.Errorf("error during database initialization: '%s'", err.Error()))
	}

	passwordFacade := password.PasswordFacade{}
	timeHelper := helper.TimeHelper{}
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

	uploadEndpoint := endpoint.NewUploadEndpoint(db, &metadataExtractor, &hasher, &jwtManager)
	router.POST("/v1/upload", uploadEndpoint.Post)
	router.OPTIONS("/v1/upload", uploadEndpoint.Options)

	fileDataEndpoint := endpoint.NewFileDataEndpoint(db, &jwtManager)
	router.GET("/v1/file_data", fileDataEndpoint.Get)

	fileEndpoint := endpoint.NewFileEndpoint(db, &jwtManager)
	router.GET("/v1/file", fileEndpoint.Get)
	router.DELETE("/v1/file", fileEndpoint.Delete)

	if len(os.Args) == 2 && os.Args[1] == "--testing" {
		restartEndpoint := endpoint.NewRestartEndpoint()
		router.POST("/v1/restart", restartEndpoint.Post)
		router.HEAD("/v1/restart", restartEndpoint.Head)
	}

	runTLSIfEnabled(router, &envGetter)
	router.Run()
}
