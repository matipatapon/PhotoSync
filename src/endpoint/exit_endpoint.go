package endpoint

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Endpoint used by ft tests to shutdown application and force it to be restarted
type ExitEndpoint struct {
	logger *log.Logger
}

func NewExitEndpoint() ExitEndpoint {
	return ExitEndpoint{logger: log.New(os.Stdout, "[ExitEndpoint]: ", log.LstdFlags)}
}

func (ee *ExitEndpoint) Post(c *gin.Context) {
	go func() {
		ee.logger.Print("Application will shutdown within 1 second")
		time.Sleep(1 * time.Second)
		ee.logger.Print("Shutting down...")
		os.Exit(0)
	}()
}
