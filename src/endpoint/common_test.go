package endpoint_test

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

var USERNAME string = "user"
var PASSWORD string = "password"
var HASH string = "HASH"

func prepareGin() (*gin.Engine, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	responseRecorder := httptest.NewRecorder()
	_, router := gin.CreateTestContext(responseRecorder)
	return router, responseRecorder
}
