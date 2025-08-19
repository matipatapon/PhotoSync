package endpoint

import (
	"os"

	"github.com/gin-gonic/gin"
)

// Endpoint used by ft tests to shutdown application and force it to be restarted
type RestartEndpoint struct {
}

func NewRestartEndpoint() RestartEndpoint {
	return RestartEndpoint{}
}

func (*RestartEndpoint) Post(c *gin.Context) {
	go func() {
		<-c.Request.Context().Done()
		os.Exit(0)
	}()
	c.Status(200)
}

func (*RestartEndpoint) Head(c *gin.Context) {
	c.Status(200)
}
