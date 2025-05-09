package endpoint

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"photosync/src/mock"

	"github.com/gin-gonic/gin"
)

var USERNAME string = "user"
var PASSWORD string = "password"
var HASH string = "HASH"
var ERROR error = errors.New("ERROR")
var INVALID_PAYLOAD = []byte("non json data")

func prepareGin() (*gin.Engine, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	responseRecorder := httptest.NewRecorder()
	_, router := gin.CreateTestContext(responseRecorder)
	return router, responseRecorder
}

func TestRegisterEndpointShouldRegisterNewUser(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectQuery("INSERT INTO users VALUES($1, $2)", [][]any{}, []any{USERNAME, HASH}, nil)
	passwordFacadeMock := mock.NewPasswordFacadeMock(t)
	passwordFacadeMock.ExpectHashPassword(PASSWORD, HASH, nil)
	sut := RegisterEndpoint{&databaseMock, &passwordFacadeMock}
	router, responseRecorder := prepareGin()
	registerData := RegisterData{USERNAME, PASSWORD}
	registerDataBytes, err := json.Marshal(registerData)
	if err != nil {
		t.Error(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(registerDataBytes)))
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 200 {
		t.Error(responseRecorder.Code)
	}

	databaseMock.AssertAllExpectionsSatisfied()
	passwordFacadeMock.AssertAllExpectionsSatisfied()
}

func TestRegisterEndpointShouldReturnErrorWhenFailedToHashPassword(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	passwordFacadeMock := mock.NewPasswordFacadeMock(t)
	passwordFacadeMock.ExpectHashPassword(PASSWORD, HASH, ERROR)
	sut := RegisterEndpoint{&databaseMock, &passwordFacadeMock}
	router, responseRecorder := prepareGin()
	registerData := RegisterData{USERNAME, PASSWORD}
	registerDataBytes, err := json.Marshal(registerData)
	if err != nil {
		t.Error(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(registerDataBytes)))
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 400 {
		t.Error(responseRecorder.Code)
	}

	databaseMock.AssertAllExpectionsSatisfied()
	passwordFacadeMock.AssertAllExpectionsSatisfied()
}

func TestRegisterEndpointShouldReturnErrorWhenRequestHasInvalidPayload(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	passwordFacadeMock := mock.NewPasswordFacadeMock(t)
	sut := RegisterEndpoint{&databaseMock, &passwordFacadeMock}
	router, responseRecorder := prepareGin()
	request := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(INVALID_PAYLOAD)))
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 400 {
		t.Error(responseRecorder.Code)
	}

	databaseMock.AssertAllExpectionsSatisfied()
	passwordFacadeMock.AssertAllExpectionsSatisfied()
}

func TestRegisterEndpointShouldReturnErrorWhenQueryFailed(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectQuery("INSERT INTO users VALUES($1, $2)", [][]any{}, []any{USERNAME, HASH}, ERROR)
	passwordFacadeMock := mock.NewPasswordFacadeMock(t)
	passwordFacadeMock.ExpectHashPassword(PASSWORD, HASH, nil)
	sut := RegisterEndpoint{&databaseMock, &passwordFacadeMock}
	router, responseRecorder := prepareGin()
	registerData := RegisterData{USERNAME, PASSWORD}
	registerDataBytes, err := json.Marshal(registerData)
	if err != nil {
		t.Error(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(registerDataBytes)))
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 400 {
		t.Error(responseRecorder.Code)
	}

	databaseMock.AssertAllExpectionsSatisfied()
	passwordFacadeMock.AssertAllExpectionsSatisfied()
}
