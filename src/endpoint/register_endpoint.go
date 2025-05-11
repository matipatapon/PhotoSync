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
	Username string `json:"username"`
	Password string `json:"password"`
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
	err := c.BindJSON(&registerData)
	if err != nil {
		logger.Print("Invalid request received!")
		c.Status(400)
		return
	}

	username := registerData.Username
	password := registerData.Password
	logger.Printf("Attempting to register '%s'", username)

	hash, err := re.passwordFacade.HashPassword(password)
	if err != nil {
		logger.Printf("Failed to hash password '%s', code 400 returned", err.Error())
		c.Status(400)
		return
	}

	err = re.db.Execute("INSERT INTO users VALUES($1, $2)", username, hash)
	if err != nil {
		logger.Printf("Execute failed with following error: '%s'", err.Error())
		c.Status(400)
		return
	}

	logger.Printf("Registered '%s'", username)
	c.Status(200)
}
