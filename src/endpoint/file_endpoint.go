package endpoint

import (
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

func (fe *FileEndpoint) Get(c *gin.Context) {
	result := fe.handleRequest(c, "SELECT file FROM files WHERE id = $1 AND user_id = $2")
	if result != nil {
		c.Writer.Write(result.QueryResult.([]byte))
		fe.logger.Printf("Returned file '%d' for user '%d'", result.FileId, result.UserId)
	}
}

func (fe *FileEndpoint) Delete(c *gin.Context) {
	result := fe.handleRequest(c, "DELETE FROM files WHERE id = $1 AND user_id = $2 RETURNING id")
	if result != nil {
		c.Status(200)
		fe.logger.Printf("Removed file '%d' for user '%d'", result.FileId, result.UserId)
	}
}

type requestResult struct {
	UserId      int64
	FileId      int64
	QueryResult any
}

func (fe *FileEndpoint) handleRequest(c *gin.Context, sql string) *requestResult {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil || id < 0 {
		c.Status(400)
		fe.logger.Print("Invalid id")
		return nil
	}

	token := c.Request.Header.Get("Authorization")
	jwt, err := fe.jm.Decode(token)
	if err != nil {
		c.Status(403)
		fe.logger.Printf("Token is invalid: '%s'", err.Error())
		return nil
	}

	rows, err := fe.db.Query(sql, id, jwt.UserId)
	if err != nil {
		c.Status(500)
		fe.logger.Printf("Query failed: '%s'", err.Error())
		return nil
	}

	if len(rows) == 0 || len(rows[0]) == 0 {
		c.Status(404)
		fe.logger.Printf("User '%d' doesn't have image with id '%d'", jwt.UserId, id)
		return nil
	}

	return &requestResult{
		FileId:      id,
		UserId:      jwt.UserId,
		QueryResult: rows[0][0],
	}
}
