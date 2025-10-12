package endpoint

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"photosync/src/database"
	"photosync/src/jwt"
	"photosync/src/metadata"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FileData struct {
	Id           string `json:"id"`
	Filename     string `json:"filename"`
	CreationDate string `json:"creation_date"`
	MIMEType     string `json:"mime_type"`
	Size         string `json:"size"`
	Thumbnail    string `json:"thumbnail"`
}

type FileDataEndpoint struct {
	db     database.IDataBase
	jm     jwt.IJwtManager
	logger *log.Logger
}

func NewFileDataEndpoint(db database.IDataBase, jm jwt.IJwtManager) FileDataEndpoint {
	return FileDataEndpoint{db: db, jm: jm, logger: log.New(os.Stdout, "[FileDataEndpoint]: ", log.LstdFlags)}
}

func (fe *FileDataEndpoint) Options(c *gin.Context) {
	c.Header("Access-Control-Allow-Headers", "Authorization")
	c.Header("Access-Control-Allow-Methods", "GET")
	c.Status(200)
}

func (fe *FileDataEndpoint) Get(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	jwt, err := fe.jm.Decode(token)
	if err != nil {
		c.Status(403)
		fe.logger.Printf("Token is invalid: '%s'", err.Error())
		return
	}

	date := c.Query("date")
	isDateOk, err := regexp.MatchString("^\\d{4}\\.\\d{2}\\.\\d{2}$", date)
	if !isDateOk {
		c.Status(400)
		fe.logger.Printf("'%s' date is invalid", date)
		return
	}
	if err != nil {
		c.Status(500)
		fe.logger.Printf("Regexp for date '%s' failed: '%s'", date, err.Error())
		return
	}

	rows, err := fe.db.Query("SELECT id, filename, TO_CHAR(creation_date, 'YYYY.MM.DD HH24:MI:SS') AS date, mime_type, size, thumbnail FROM files WHERE user_id = $1 AND TO_CHAR(creation_date, 'YYYY.MM.DD') ILIKE $2 || '%' ORDER BY id DESC, creation_date DESC", jwt.UserId, date)
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

		var thumbnail string = ""
		if row[5] != nil {
			thumbnail = base64.StdEncoding.EncodeToString(row[5].([]byte))
		} else {
			file, err := fe.db.Query("SELECT file FROM files WHERE id = $1", id)
			if err != nil {
				c.Status(500)
				fe.logger.Printf("Image query failed: '%s'", err.Error())
				return
			}
			if len(file) == 0 || len(file[0]) == 0 {
				c.Status(500)
				fe.logger.Printf("No image with id {%d}, probably got removed", id)
				return
			}
			thumbnail = base64.StdEncoding.EncodeToString(file[0][0].([]byte))
		}

		fileData = append(fileData, FileData{
			Id:           strconv.FormatInt(id, 10),
			Filename:     filename,
			CreationDate: creation_date,
			MIMEType:     mime_type,
			Size:         strconv.FormatInt(size, 10),
			Thumbnail:    thumbnail,
		})
	}

	bytes, _ := json.Marshal(fileData)
	c.Writer.Write(bytes)
}
