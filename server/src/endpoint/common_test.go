package endpoint_test

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

var TOKEN_STRING string = "TOKEN_STRING"
var EXPIRATION_TIME int64 = 102

var USERNAME string = "user"
var PASSWORD string = "password"
var USER_ID int64 = 666

var FILE []byte = []byte("FILE_CONTENT")
var FILENAME string = "dog.jpg"
var HASH string = "HASH"
var FILE_ID int64 = 2137

var INVALID_MODIFICATION_DATE string = "2025:04:12|15-32-13"
var MODIFICATION_DATE string = "2025.08.03 15:24:13"
var CREATION_DATE string = "2024.07.01 12:31:32"
var PARAM_DATE string = "2025.05.16"

func prepareGin() (*gin.Engine, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	responseRecorder := httptest.NewRecorder()
	_, router := gin.CreateTestContext(responseRecorder)
	return router, responseRecorder
}
