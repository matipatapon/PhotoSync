package endpoint

import (
	"errors"
	"log"
	"os"
	"photosync/src/database"
	"photosync/src/jwt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FileEndpoint struct {
	db     database.IDataBase
	jm     jwt.IJwtManager
	logger *log.Logger
}

func NewFileEndpoint(db database.IDataBase, jm jwt.IJwtManager) FileEndpoint {
	return FileEndpoint{db: db, jm: jm, logger: log.New(os.Stdout, "[FileEndpoint]: ", log.LstdFlags)}
}

// TODO docs & fties
func (fe *FileEndpoint) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil || id < 0 {
		c.Status(400)
		fe.logger.Print("Invalid id")
		return
	}

	token := c.Request.Header.Get("Authorization")
	jwt, err := fe.jm.Decode(token)
	if err != nil {
		c.Status(403)
		fe.logger.Printf("Token is invalid: '%s'", err.Error())
		return
	}

	rows, err := fe.db.Query("SELECT file FROM files WHERE id = $1 AND user_id = $2", id, jwt.UserId)
	if err != nil {
		c.Status(500)
		fe.logger.Printf("Query failed: '%s'", err.Error())
		return
	}

	if len(rows) == 0 || len(rows[0]) == 0 {
		c.Status(404)
		fe.logger.Printf("User '%d' doesn't have image with id '%d'", jwt.UserId, id)
		return
	}

	c.Writer.Write(rows[0][0].([]byte))
}

func (fe *FileEndpoint) Delete(c *gin.Context) {
	panic(errors.New("not implemented"))
}
