package endpoint

import (
	"os"

	"github.com/gin-gonic/gin"
)

// Endpoint used by ft tests to shutdown application and force it to be restarted
type RestartEndpoint struct {
	restart bool
}

func NewRestartEndpoint() RestartEndpoint {
	return RestartEndpoint{restart: false}
}

func (re *RestartEndpoint) Post(c *gin.Context) {
	re.restart = true
	c.Status(200)
}

func (re *RestartEndpoint) Head(c *gin.Context) {
	if re.restart {
		os.Exit(0)
	}
	c.Status(200)
}
