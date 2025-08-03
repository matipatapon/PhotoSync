package endpoint_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"photosync/src/endpoint"
	"photosync/src/jwt"
	"photosync/src/mock"
	"testing"
)

var LOGIN_SQL string = "SELECT id, password FROM users WHERE username = $1"
var ONE_DAY int64 = 60 * 60 * 24

type FakeReader struct {
}

func (FakeReader) Close() error {
	return nil
}

func (FakeReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

func TestLoginEndpointShouldSendJwtWhenCredentialsAreCorrect(t *testing.T) {
	dbMock := mock.NewDatabaseMock(t)
	pfMock := mock.NewPasswordFacadeMock(t)
	jmMock := mock.NewJwtManagerMock(t)
	thMock := mock.NewTimeHelperMock(t)
	sut := endpoint.NewLoginEndpoint(&dbMock, &pfMock, &jmMock, &thMock)
	loginData := endpoint.LoginData{Username: USERNAME, Password: PASSWORD}
	loginDataBytes, err := json.Marshal(loginData)
	if err != nil {
		t.FailNow()
	}
	jwtPayload := jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: int64(EXPIRATION_TIME)}

	pfMock.ExpectMatchHashToPassword(HASH, PASSWORD, true)
	dbMock.ExpectQuery(LOGIN_SQL, [][]any{{USER_ID, HASH}}, []any{USERNAME}, nil)
	jmMock.ExpectCreate(jwtPayload, TOKEN_STRING, nil)
	thMock.ExpectTimeIn(ONE_DAY, int64(EXPIRATION_TIME))

	router, responseRecorder := prepareGin()
	request := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(loginDataBytes)))
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 200 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.String() != TOKEN_STRING {
		t.Errorf("wrong token '%s'", responseRecorder.Body.String())
	}

	pfMock.AssertAllExpectionsSatisfied()
	dbMock.AssertAllExpectionsSatisfied()
	jmMock.AssertAllExpectionsSatisfied()
	thMock.AssertAllExpectionsSatisfied()
}

func TestLoginEndpointShouldReturn500WhenFailedToReadRequestBody(t *testing.T) {
	dbMock := mock.NewDatabaseMock(t)
	pfMock := mock.NewPasswordFacadeMock(t)
	jmMock := mock.NewJwtManagerMock(t)
	thMock := mock.NewTimeHelperMock(t)
	sut := endpoint.NewLoginEndpoint(&dbMock, &pfMock, &jmMock, &thMock)
	loginData := endpoint.LoginData{Username: USERNAME, Password: PASSWORD}
	loginDataBytes, err := json.Marshal(loginData)
	if err != nil {
		t.FailNow()
	}
	router, responseRecorder := prepareGin()
	request := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(loginDataBytes)))
	request.Body = FakeReader{}
	router.POST("/", sut.Post)

	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 500 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "" {
		t.Error("request body should be empty")
	}

	pfMock.AssertAllExpectionsSatisfied()
	dbMock.AssertAllExpectionsSatisfied()
	jmMock.AssertAllExpectionsSatisfied()
	thMock.AssertAllExpectionsSatisfied()
}

func TestLoginEndpointShouldReturn400WhenInvalidRequestBody(t *testing.T) {
	dbMock := mock.NewDatabaseMock(t)
	pfMock := mock.NewPasswordFacadeMock(t)
	jmMock := mock.NewJwtManagerMock(t)
	thMock := mock.NewTimeHelperMock(t)
	sut := endpoint.NewLoginEndpoint(&dbMock, &pfMock, &jmMock, &thMock)

	router, responseRecorder := prepareGin()
	request := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader([]byte{})))
	router.POST("/", sut.Post)

	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 400 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "" {
		t.Error("request body should be empty")
	}

	pfMock.AssertAllExpectionsSatisfied()
	dbMock.AssertAllExpectionsSatisfied()
	jmMock.AssertAllExpectionsSatisfied()
	thMock.AssertAllExpectionsSatisfied()
}

func TestLoginEndpointShouldReturn400WhenUserDoesNotExistNoRowsReturned(t *testing.T) {
	dbMock := mock.NewDatabaseMock(t)
	pfMock := mock.NewPasswordFacadeMock(t)
	jmMock := mock.NewJwtManagerMock(t)
	thMock := mock.NewTimeHelperMock(t)
	sut := endpoint.NewLoginEndpoint(&dbMock, &pfMock, &jmMock, &thMock)
	loginData := endpoint.LoginData{Username: USERNAME, Password: PASSWORD}
	loginDataBytes, err := json.Marshal(loginData)
	if err != nil {
		t.FailNow()
	}

	dbMock.ExpectQuery(LOGIN_SQL, [][]any{}, []any{USERNAME}, nil)

	router, responseRecorder := prepareGin()
	request := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(loginDataBytes)))
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 400 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "" {
		t.Error("request body should be empty")
	}

	pfMock.AssertAllExpectionsSatisfied()
	dbMock.AssertAllExpectionsSatisfied()
	jmMock.AssertAllExpectionsSatisfied()
	thMock.AssertAllExpectionsSatisfied()
}

func TestLoginEndpointShouldReturn400WhenUserDoesNotExistEmptyRowReturned(t *testing.T) {
	dbMock := mock.NewDatabaseMock(t)
	pfMock := mock.NewPasswordFacadeMock(t)
	jmMock := mock.NewJwtManagerMock(t)
	thMock := mock.NewTimeHelperMock(t)
	sut := endpoint.NewLoginEndpoint(&dbMock, &pfMock, &jmMock, &thMock)
	loginData := endpoint.LoginData{Username: USERNAME, Password: PASSWORD}
	loginDataBytes, err := json.Marshal(loginData)
	if err != nil {
		t.FailNow()
	}

	dbMock.ExpectQuery(LOGIN_SQL, [][]any{{}}, []any{USERNAME}, nil)

	router, responseRecorder := prepareGin()
	request := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(loginDataBytes)))
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 400 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "" {
		t.Error("request body should be empty")
	}

	pfMock.AssertAllExpectionsSatisfied()
	dbMock.AssertAllExpectionsSatisfied()
	jmMock.AssertAllExpectionsSatisfied()
	thMock.AssertAllExpectionsSatisfied()
}

func TestLoginEndpointShouldReturn400WhenPasswordIsInvalid(t *testing.T) {
	dbMock := mock.NewDatabaseMock(t)
	pfMock := mock.NewPasswordFacadeMock(t)
	jmMock := mock.NewJwtManagerMock(t)
	thMock := mock.NewTimeHelperMock(t)
	sut := endpoint.NewLoginEndpoint(&dbMock, &pfMock, &jmMock, &thMock)
	loginData := endpoint.LoginData{Username: USERNAME, Password: PASSWORD}
	loginDataBytes, err := json.Marshal(loginData)
	if err != nil {
		t.FailNow()
	}

	pfMock.ExpectMatchHashToPassword(HASH, PASSWORD, false)
	dbMock.ExpectQuery(LOGIN_SQL, [][]any{{USER_ID, HASH}}, []any{USERNAME}, nil)

	router, responseRecorder := prepareGin()
	request := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(loginDataBytes)))
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 400 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "" {
		t.Error("request body should be empty")
	}

	pfMock.AssertAllExpectionsSatisfied()
	dbMock.AssertAllExpectionsSatisfied()
	jmMock.AssertAllExpectionsSatisfied()
	thMock.AssertAllExpectionsSatisfied()
}

func TestLoginEndpointShouldReturn500WhenFailedToCreateToken(t *testing.T) {
	dbMock := mock.NewDatabaseMock(t)
	pfMock := mock.NewPasswordFacadeMock(t)
	jmMock := mock.NewJwtManagerMock(t)
	thMock := mock.NewTimeHelperMock(t)
	sut := endpoint.NewLoginEndpoint(&dbMock, &pfMock, &jmMock, &thMock)
	loginData := endpoint.LoginData{Username: USERNAME, Password: PASSWORD}
	loginDataBytes, err := json.Marshal(loginData)
	if err != nil {
		t.FailNow()
	}
	jwtPayload := jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: int64(EXPIRATION_TIME)}

	pfMock.ExpectMatchHashToPassword(HASH, PASSWORD, true)
	dbMock.ExpectQuery(LOGIN_SQL, [][]any{{USER_ID, HASH}}, []any{USERNAME}, nil)
	jmMock.ExpectCreate(jwtPayload, TOKEN_STRING, errors.New("error"))
	thMock.ExpectTimeIn(ONE_DAY, int64(EXPIRATION_TIME))

	router, responseRecorder := prepareGin()
	request := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(loginDataBytes)))
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 500 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "" {
		t.Error("request body should be empty")
	}

	pfMock.AssertAllExpectionsSatisfied()
	dbMock.AssertAllExpectionsSatisfied()
	jmMock.AssertAllExpectionsSatisfied()
	thMock.AssertAllExpectionsSatisfied()
}
