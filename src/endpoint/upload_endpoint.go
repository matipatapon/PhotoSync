package endpoint

import (
	"log"
	"photosync/src/database"
	"photosync/src/jwt"

	"github.com/gin-gonic/gin"
)

type UploadEndpoint struct {
	db     database.IDataBase
	jm     jwt.IJwtManager
	logger *log.Logger
}

func (le *UploadEndpoint) Post(c *gin.Context) {
}
