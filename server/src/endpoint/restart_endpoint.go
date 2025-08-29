package endpoint

import (
	"log"
	"os"
	"photosync/src/database"

	"github.com/gin-gonic/gin"
)

type RestartEndpoint struct {
	db     database.IDataBase
	logger *log.Logger
}

func NewRestartEndpoint(db database.IDataBase) RestartEndpoint {
	return RestartEndpoint{db: db, logger: log.New(os.Stdout, "[RestartEndpoint]: ", log.LstdFlags)}
}

func (re *RestartEndpoint) Post(c *gin.Context) {
	err := re.db.DropDb()
	if err != nil {
		re.logger.Printf("Failed to drop db: '%s'", err.Error())
		c.Status(500)
		return
	}

	err = re.db.InitDb()
	if err != nil {
		re.logger.Printf("Failed to recreata db: '%s'", err.Error())
		c.Status(500)
		return
	}

	c.Status(200)
}
