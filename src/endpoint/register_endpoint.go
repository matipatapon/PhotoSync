package endpoint

import (
	"log"
	"os"
	"photosync/src/database"
	"photosync/src/password"

	"github.com/gin-gonic/gin"
)

type RegisterData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterEndpoint struct {
	db             database.IDataBase
	passwordFacade password.IPasswordFacade
	logger         *log.Logger
}

func NewRegisterEndpoint(db database.IDataBase, passwordFacade password.IPasswordFacade) RegisterEndpoint {
	return RegisterEndpoint{db, passwordFacade, log.New(os.Stdout, "[RegisterEndpoint]: ", log.LstdFlags)}
}

func (re *RegisterEndpoint) Post(c *gin.Context) {
	var registerData RegisterData
	err := c.BindJSON(&registerData)
	if err != nil {
		re.logger.Print("Invalid request received!")
		c.Status(400)
		return
	}

	username := registerData.Username
	if username == "" {
		re.logger.Printf("No username given")
		c.Status(400)
		return
	}
	re.logger.Printf("Attempting to register '%s'", username)

	password := registerData.Password
	if password == "" {
		re.logger.Printf("No password given")
		c.Status(400)
		return
	}

	hash, err := re.passwordFacade.HashPassword(password)
	if err != nil {
		re.logger.Printf("Failed to hash password '%s', code 500 returned", err.Error())
		c.Status(500)
		return
	}

	err = re.db.Execute("INSERT INTO users(username, password) VALUES($1, $2)", username, hash)
	if err != nil {
		re.logger.Printf("Execute failed with following error: '%s'", err.Error())
		c.Status(400)
		return
	}

	re.logger.Printf("Registered '%s'", username)
	c.Status(200)
}
