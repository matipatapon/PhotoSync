package endpoint

import (
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type UploadEndpoint struct {
	logger *log.Logger
}

func NewUploadEndpoint() UploadEndpoint {
	return UploadEndpoint{logger: log.New(os.Stdout, "[UploadEndpoint]: ", log.LstdFlags)}
}

func (ue *UploadEndpoint) Post(c *gin.Context) {
	reader, _ := c.Request.MultipartReader()
	p, _ := reader.NextPart()
	filenameBytes, _ := io.ReadAll(p)
	filename := string(filenameBytes)
	p, _ = reader.NextPart()
	modificationDateBytes, _ := io.ReadAll(p)
	modification_date := string(modificationDateBytes)
	// p, _ = reader.NextPart()
	// file, _ := io.ReadAll(p)

	ue.logger.Printf("Filename: '%s'", filename)
	ue.logger.Printf("Modification date: '%s'", modification_date)
}
