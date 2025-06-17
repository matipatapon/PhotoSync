package endpoint_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"photosync/src/endpoint"
	"photosync/src/mock"
)

var ERROR error = errors.New("ERROR")
var INVALID_PAYLOAD = []byte("non json data")

func TestRegisterEndpointShouldRegisterNewUser(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectExecute("INSERT INTO users(username, password) VALUES($1, $2)", []any{USERNAME, HASH}, nil)
	passwordFacadeMock := mock.NewPasswordFacadeMock(t)
	passwordFacadeMock.ExpectHashPassword(PASSWORD, HASH, nil)
	sut := endpoint.NewRegisterEndpoint(&databaseMock, &passwordFacadeMock)
	router, responseRecorder := prepareGin()
	registerData := endpoint.RegisterData{USERNAME, PASSWORD}
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
	sut := endpoint.NewRegisterEndpoint(&databaseMock, &passwordFacadeMock)
	router, responseRecorder := prepareGin()
	registerData := endpoint.RegisterData{USERNAME, PASSWORD}
	registerDataBytes, err := json.Marshal(registerData)
	if err != nil {
		t.Error(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(registerDataBytes)))
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 500 {
		t.Error(responseRecorder.Code)
	}

	databaseMock.AssertAllExpectionsSatisfied()
	passwordFacadeMock.AssertAllExpectionsSatisfied()
}

func TestRegisterEndpointShouldReturnErrorWhenRequestHasInvalidPayload(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	passwordFacadeMock := mock.NewPasswordFacadeMock(t)
	sut := endpoint.NewRegisterEndpoint(&databaseMock, &passwordFacadeMock)
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
	databaseMock.ExpectExecute("INSERT INTO users(username, password) VALUES($1, $2)", []any{USERNAME, HASH}, ERROR)
	passwordFacadeMock := mock.NewPasswordFacadeMock(t)
	passwordFacadeMock.ExpectHashPassword(PASSWORD, HASH, nil)
	sut := endpoint.NewRegisterEndpoint(&databaseMock, &passwordFacadeMock)
	router, responseRecorder := prepareGin()
	registerData := endpoint.RegisterData{USERNAME, PASSWORD}
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

func TestRegisterEndpointShouldNotRegisterWhenNoUsernameGiven(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	passwordFacadeMock := mock.NewPasswordFacadeMock(t)
	sut := endpoint.NewRegisterEndpoint(&databaseMock, &passwordFacadeMock)
	router, responseRecorder := prepareGin()

	type InvalidRegisterData struct {
		Password string `json:"password"`
	}
	registerData := InvalidRegisterData{PASSWORD}

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

func TestRegisterEndpointShouldNotRegisterWhenNoPasswordGiven(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	passwordFacadeMock := mock.NewPasswordFacadeMock(t)
	sut := endpoint.NewRegisterEndpoint(&databaseMock, &passwordFacadeMock)
	router, responseRecorder := prepareGin()

	type InvalidRegisterData struct {
		Username string `json:"username"`
	}
	registerData := InvalidRegisterData{USERNAME}

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
