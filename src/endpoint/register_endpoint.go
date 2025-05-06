package endpoint

import (
	"log"
	"os"
	"photosync/src/database"
	"photosync/src/password"

	"github.com/gin-gonic/gin"
)

var logger *log.Logger = log.New(os.Stdout, "[RegisterEndpoint]: ", log.LstdFlags)

type RegisterData struct {
	Username string `json: username`
	Password string `json: password`
}

type RegisterEndpoint struct {
	db             database.IDataBase
	passwordFacade password.IPasswordFacade
}

func NewRegisterEndpoint(db database.IDataBase, passwordFacade password.IPasswordFacade) RegisterEndpoint {
	return RegisterEndpoint{db, passwordFacade}
}

func (re *RegisterEndpoint) Post(c *gin.Context) {
	var registerData RegisterData
	c.BindJSON(&registerData)
	username := registerData.Username
	password := registerData.Password

	logger.Printf("Attempting to register '%s'", username)

	hash, err := re.passwordFacade.HashPassword(password)
	if err != nil {
		logger.Printf("Failed to hash password '%s', code 400 returned", err.Error())
		c.Status(400)
		return
	}
	re.db.Query("INSERT INTO users VALUES($1, $2)", username, hash)
	c.Status(200)
}
