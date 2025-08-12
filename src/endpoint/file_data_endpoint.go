package endpoint

import (
	"encoding/json"
	"log"
	"os"
	"photosync/src/database"
	"photosync/src/jwt"
	"photosync/src/metadata"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FileData struct {
	Id           string `json:"id"`
	Filename     string `json:"filename"`
	CreationDate string `json:"creation_date"`
	MIMEType     string `json:"mime_type"`
	Size         string `json:"size"`
}

type FileDataEndpoint struct {
	db     database.IDataBase
	jm     jwt.IJwtManager
	logger *log.Logger
}

func NewFileDataEndpoint(db database.IDataBase, jm jwt.IJwtManager) FileDataEndpoint {
	return FileDataEndpoint{db: db, jm: jm, logger: log.New(os.Stdout, "[FileDataEndpoint]: ", log.LstdFlags)}
}

func (fe *FileDataEndpoint) Get(c *gin.Context) {
	offset, err := strconv.ParseInt(c.Query("offset"), 10, 64)
	if err != nil || offset < 0 {
		c.Status(400)
		fe.logger.Print("Invalid offset")
		return
	}

	count, err := strconv.ParseInt(c.Query("count"), 10, 64)
	if err != nil || count < 0 {
		c.Status(400)
		fe.logger.Print("Invalid count")
		return
	}

	token := c.Request.Header.Get("Authorization")
	jwt, err := fe.jm.Decode(token)
	if err != nil {
		c.Status(403)
		fe.logger.Printf("Token is invalid: '%s'", err.Error())
		return
	}

	rows, err := fe.db.Query("SELECT id, filename, TO_CHAR(creation_date, 'YYYY.MM.DD HH24:MI:SS'), mime_type, size FROM files WHERE user_id = $1 ORDER BY creation_date DESC LIMIT $2 OFFSET $3", jwt.UserId, count, offset)
	if err != nil {
		c.Status(500)
		fe.logger.Printf("Query failed: '%s'", err.Error())
		return
	}

	fileData := []FileData{}
	for _, row := range rows {
		id := row[0].(int64)
		filename := row[1].(string)
		creation_date := row[2].(string)
		mime_type_raw := row[3].(int16)
		mime_type := metadata.MIMETypeToString(metadata.MIMEType(mime_type_raw))
		size := row[4].(int64)

		fileData = append(fileData, FileData{
			Id:           strconv.FormatInt(id, 10),
			Filename:     filename,
			CreationDate: creation_date,
			MIMEType:     mime_type,
			Size:         strconv.FormatInt(size, 10),
		})
	}

	bytes, _ := json.Marshal(fileData)
	c.Writer.Write(bytes)
}
