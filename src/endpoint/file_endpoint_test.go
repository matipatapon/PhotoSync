package endpoint_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"photosync/src/endpoint"
	"photosync/src/jwt"
	"photosync/src/mock"
	"reflect"
	"testing"
)

var FILE_ENDPOINT_SQL string = "SELECT file FROM files WHERE id = $1 AND user_id = $2"

func prepareGetFileRequest(id *int64) *http.Request {
	var query string
	if id != nil {
		query = fmt.Sprintf("/?id=%d", *id)
	} else {
		query = "/"
	}
	request := httptest.NewRequest(http.MethodGet, query, io.NopCloser(bytes.NewBuffer([]byte{})))
	request.Header.Set("Authorization", TOKEN_STRING)
	return request
}

func TestFileEndpointShouldReturn404WhenImageNotExists(t *testing.T) {
	requests := [][][]any{
		{},
		{{}},
	}
	for _, request := range requests {
		databaseMock := mock.NewDatabaseMock(t)
		databaseMock.ExpectQuery(FILE_ENDPOINT_SQL, request, []any{FILE_ID, USER_ID}, nil)
		defer databaseMock.AssertAllExpectionsSatisfied()

		jwtManagerMock := mock.NewJwtManagerMock(t)
		jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID}, nil)
		defer jwtManagerMock.AssertAllExpectionsSatisfied()

		sut := endpoint.NewFileEndpoint(&databaseMock, &jwtManagerMock)

		router, responseRecorder := prepareGin()
		router.GET("/", sut.Get)
		router.ServeHTTP(responseRecorder, prepareGetFileRequest(&FILE_ID))

		if responseRecorder.Code != 404 {
			t.Error(responseRecorder.Code)
		}
		if responseRecorder.Body.Len() != 0 {
			t.Error("Expected empty response")
		}
	}
}

func TestFileEndpointShouldReturn500WhenQueryFailed(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectQuery(FILE_ENDPOINT_SQL, [][]any{{FILE}}, []any{FILE_ID, USER_ID}, errors.New("query error"))
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewFileEndpoint(&databaseMock, &jwtManagerMock)

	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)
	router.ServeHTTP(responseRecorder, prepareGetFileRequest(&FILE_ID))

	if responseRecorder.Code != 500 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.Len() != 0 {
		t.Error("Expected empty response")
	}
}

func TestFileEndpointShouldReturn403WhenTokenIsInvalid(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID}, errors.New("invalid token"))
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewFileEndpoint(&databaseMock, &jwtManagerMock)

	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)
	router.ServeHTTP(responseRecorder, prepareGetFileRequest(&FILE_ID))

	if responseRecorder.Code != 403 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.Len() != 0 {
		t.Error("Expected empty response")
	}
}

func TestFileEndpointShouldReturnFile(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectQuery(FILE_ENDPOINT_SQL, [][]any{{FILE}}, []any{FILE_ID, USER_ID}, nil)
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewFileEndpoint(&databaseMock, &jwtManagerMock)

	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)
	router.ServeHTTP(responseRecorder, prepareGetFileRequest(&FILE_ID))

	if responseRecorder.Code != 200 {
		t.Error(responseRecorder.Code)
	}
	if !reflect.DeepEqual(responseRecorder.Body.Bytes(), FILE) {
		t.Error("Unexpected body")
	}
}

func TestFileEndpointShouldReturn400WhenIdNotSpecified(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewFileEndpoint(&databaseMock, &jwtManagerMock)

	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)
	router.ServeHTTP(responseRecorder, prepareGetFileRequest(nil))

	if responseRecorder.Code != 400 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.Len() != 0 {
		t.Error("Expected empty response")
	}
}
