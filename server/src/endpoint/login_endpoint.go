package endpoint

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"photosync/src/database"
	"photosync/src/helper"
	"photosync/src/jwt"
	"photosync/src/password"

	"github.com/gin-gonic/gin"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginEndpoint struct {
	db     database.IDataBase
	pf     password.IPasswordFacade
	jm     jwt.IJwtManager
	th     helper.ITimeHelper
	logger *log.Logger
}

func NewLoginEndpoint(db database.IDataBase, pf password.IPasswordFacade, jm jwt.IJwtManager, th helper.ITimeHelper) LoginEndpoint {
	return LoginEndpoint{db: db, pf: pf, jm: jm, th: th, logger: log.New(os.Stdout, "[LoginEndpoint]: ", log.LstdFlags)}
}

func (le *LoginEndpoint) Post(c *gin.Context) {
	var loginData LoginData
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		le.logger.Printf("Failed to read request body: '%s'", err.Error())
		c.Status(400)
		return
	}

	err = json.Unmarshal(bodyBytes, &loginData)
	if err != nil {
		le.logger.Printf("Failed to retrieve login data from request body: '%s'", err.Error())
		c.Status(400)
		return
	}

	result, _ := le.db.Query("SELECT id, password FROM users WHERE username = $1", loginData.Username)
	if len(result) == 0 || len(result[0]) == 0 {
		le.logger.Printf("User '%s' does not exist in db", loginData.Username)
		c.Status(401)
		return
	}
	userId := result[0][0].(int64)
	hashedPassword := result[0][1].(string)

	if !le.pf.MatchHashToPassword(hashedPassword, loginData.Password) {
		le.logger.Printf("Password mismatch for '%s'", loginData.Username)
		c.Status(401)
		return
	}

	tokenString, err := le.jm.Create(jwt.JwtPayload{UserId: userId, Username: loginData.Username, ExpirationTime: le.th.TimeIn(60 * 60 * 24)})
	if err != nil {
		le.logger.Printf("Failed to create token")
		c.Status(500)
		return
	}

	le.logger.Printf("Successfully authenticated: '%s'", loginData.Username)
	c.String(200, tokenString)
}
