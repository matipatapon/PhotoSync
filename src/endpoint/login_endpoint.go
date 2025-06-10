package endpoint

import (
	"photosync/src/database"
	"photosync/src/jwt"
	"photosync/src/password"

	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	Username string
}

type LoginEndpoint struct {
	db database.IDataBase
	pf password.IPasswordFacade
	jm jwt.IJwtManager
}

// TODO
func (le *LoginEndpoint) Post(c *gin.Context) {
}
