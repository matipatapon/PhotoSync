package endpoint

import (
	"errors"
	"io"
	"log"
	"os"
	"photosync/src/database"
	"photosync/src/helper"
	"photosync/src/jwt"
	"photosync/src/metadata"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UploadEndpoint struct {
	db     database.IDataBase
	me     metadata.IMetadataExtractor
	h      helper.IHasher
	jm     jwt.IJwtManager
	logger *log.Logger
}

func NewUploadEndpoint(db database.IDataBase, me metadata.IMetadataExtractor, h helper.IHasher, jm jwt.IJwtManager) UploadEndpoint {
	return UploadEndpoint{
		db:     db,
		me:     me,
		h:      h,
		jm:     jm,
		logger: log.New(os.Stdout, "[UploadEndpoint]: ", log.LstdFlags),
	}
}

func (ue *UploadEndpoint) Options(c *gin.Context) {
	c.Header("Access-Control-Allow-Headers", "Authorization")
	c.Header("Access-Control-Allow-Methods", "POST")
	c.Status(200)
}

func (ue *UploadEndpoint) Post(c *gin.Context) {
	jwt, err := ue.authorize(c)
	if err != nil {
		ue.logger.Print("Token is invalid")
		c.Status(403)
		return
	}
	ue.logger.Printf("User '%s' authorized for upload", jwt.Username)

	filename, modificationDate, file, err := ue.tryToGetDataFromBody(c)
	if err != nil {
		ue.logger.Printf("Failed to extract data from request body: '%s'", err.Error())
		c.Status(400)
		return
	}

	hash, err := ue.h.Hash(file)
	if err != nil {
		ue.logger.Printf("Failed to hash a file: '%s'", err.Error())
		c.Status(500)
		return
	}

	meta := ue.me.Extract(file)
	if meta.MIMEType == metadata.UNKNOWN {
		ue.logger.Printf("Unknown file type")
		c.Status(401)
		return
	}
	if meta.CreationDate != nil {
		modificationDate = meta.CreationDate
	}

	result, err := ue.db.Query(
		"INSERT INTO files(user_id, creation_date, filename, mime_type, file, hash, size) VALUES($1, TO_TIMESTAMP($2, 'YYYY.MM.DD HH24:MI:SS'), $3, $4, $5, $6, $7) RETURNING id",
		jwt.UserId,
		modificationDate.ToString(),
		filename,
		meta.MIMEType,
		file,
		hash,
		len(file),
	)
	if err != nil {
		ue.logger.Printf("Query error: '%s'", err.Error())
		c.Status(500)
		return
	}
	if len(result) == 0 || len(result[0]) == 0 {
		ue.logger.Print("File already exists")
		c.Status(402)
		return
	}

	ue.logger.Print("Sucessfully saved a file")
	c.String(200, strconv.FormatInt(result[0][0].(int64), 10))
}

func (ue *UploadEndpoint) authorize(c *gin.Context) (jwt.JwtPayload, error) {
	return ue.jm.Decode(c.Request.Header.Get("Authorization"))
}

func (ue *UploadEndpoint) tryToGetDataFromBody(c *gin.Context) (string, *metadata.Date, []byte, error) {
	var filename string
	var modificationDate *metadata.Date
	var file []byte

	reader, err := c.Request.MultipartReader()
	if err != nil {
		ue.logger.Printf("Failed to read request: '%s'", err.Error())
		return "", nil, []byte{}, err
	}

	for {
		p, err := reader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			ue.logger.Printf("Error occured during reading parts: '%s'", err.Error()) // untested, shouldn't happen
			return "", nil, []byte{}, err
		}

		bytes, err := io.ReadAll(p)
		if err != nil {
			ue.logger.Printf("Error occured during reading part: '%s'", err.Error()) // untested, shouldn't happen
			return "", nil, []byte{}, err
		}

		switch p.FormName() {
		case "filename":
			filename = string(bytes)
		case "modification_date":
			date, err := metadata.NewDate(string(bytes))
			if err != nil {
				ue.logger.Printf("Invalid modification date: '%s'", err.Error())
				return "", nil, []byte{}, err
			}
			modificationDate = &date
		case "file":
			file = bytes
		}
	}

	if filename == "" {
		ue.logger.Print("Filename is missing")
		return "", nil, []byte{}, errors.New("filename is missing")
	}
	if modificationDate == nil {
		ue.logger.Print("Modification date is missing")
		return "", nil, []byte{}, errors.New("modification date is missing")
	}
	if len(file) == 0 {
		ue.logger.Print("File is missing")
		return "", nil, []byte{}, errors.New("file is missing")
	}

	return filename, modificationDate, file, nil
}
