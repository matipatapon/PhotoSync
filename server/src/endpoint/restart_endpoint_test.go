package endpoint_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"photosync/src/endpoint"
	"photosync/src/mock"
	"testing"
)

func TestRestartEndpointShouldReinitializeDb(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectDropDb(nil)
	databaseMock.ExpectInitDb(nil)
	defer databaseMock.AssertAllExpectionsSatisfied()
	sut := endpoint.NewRestartEndpoint(&databaseMock)

	router, responseRecorder := prepareGin()
	request := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader([]byte{})))
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 200 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.Len() != 0 {
		t.Error("expected empty response body")
	}
}

func TestRestartEndpointShouldReturn500WhenFailedToDropDb(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectDropDb(errors.New("failed to drop db"))
	defer databaseMock.AssertAllExpectionsSatisfied()
	sut := endpoint.NewRestartEndpoint(&databaseMock)

	router, responseRecorder := prepareGin()
	request := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader([]byte{})))
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 500 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.Len() != 0 {
		t.Error("expected empty response body")
	}
}

func TestRestartEndpointShouldReturn500WhenFailedToInitDb(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectDropDb(nil)
	databaseMock.ExpectInitDb(errors.New("failed to init db"))
	defer databaseMock.AssertAllExpectionsSatisfied()
	sut := endpoint.NewRestartEndpoint(&databaseMock)

	router, responseRecorder := prepareGin()
	request := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader([]byte{})))
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 500 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.Len() != 0 {
		t.Error("expected empty response body")
	}
}
