package endpoint

import (
	"os"

	"github.com/gin-gonic/gin"
)

// Endpoint used by ft tests to shutdown application and force it to be restarted
type ExitEndpoint struct {
}

func NewExitEndpoint() ExitEndpoint {
	return ExitEndpoint{}
}

func (ee *ExitEndpoint) Post(c *gin.Context) {
	go func() {
		<-c.Request.Context().Done()
		os.Exit(0)
	}()
	c.Status(200)
}

func (ee *ExitEndpoint) Head(c *gin.Context) {
	c.Status(200)
}
