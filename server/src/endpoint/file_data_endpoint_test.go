package endpoint_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"photosync/src/endpoint"
	"photosync/src/jwt"
	"photosync/src/metadata"
	"photosync/src/mock"
	"reflect"
	"strconv"
	"testing"
)

var FILE_DATA_QUERY string = "SELECT id, filename, TO_CHAR(creation_date, 'YYYY.MM.DD HH24:MI:SS'), mime_type, size FROM files WHERE user_id = $1 ORDER BY creation_date DESC, id DESC LIMIT $2 OFFSET $3"
var OFFSET int64 = 10
var COUNT int64 = 15
var NEGATIVE_OFFSET int64 = -10
var NEGATIVE_COUNT int64 = -15

var FILE_DATA []endpoint.FileData = []endpoint.FileData{
	{
		Id:           "13",
		Filename:     "B2.jpg",
		CreationDate: "2023.01.04 12:30:00",
		MIMEType:     "image/jpeg",
		Size:         "1024",
	},
	{
		Id:           "17",
		Filename:     "9S.jpg",
		CreationDate: "2022.04.06 14:33:12",
		MIMEType:     "image/jpeg",
		Size:         "5000",
	},
}

func fileDataToRequestArgs(fileData []endpoint.FileData) [][]any {
	result := [][]any{}
	for _, fd := range fileData {
		id, _ := strconv.ParseInt(fd.Id, 10, 64)
		size, _ := strconv.ParseInt(fd.Size, 10, 64)
		result = append(result, []any{
			id,
			fd.Filename,
			fd.CreationDate,
			int16(metadata.StringToMIMEType(fd.MIMEType)),
			size,
		})
	}
	return result
}

var FILE_1_ID int64 = 13
var FILE_1_FILENAME string = "Joseph_Joestar.jpg"
var FILE_1_CREATION_DATE string = "2023.01.04 12:30:00"
var FILE_1_MIME_TYPE int16 = int16(metadata.JPG)
var FILE_1_SIZE int64 = 1024

var FILE_2_ID int64 = 17
var FILE_2_FILENAME string = "A2.jpg"
var FILE_2_CREATION_DATE string = "2022.01.04 12:30:00"
var FILE_2_MIME_TYPE int16 = int16(metadata.JPG)
var FILE_2_SIZE int64 = 5000

func prepareRequest(offset *int64, count *int64) *http.Request {
	var query string
	if offset == nil {
		if count == nil {
			query = "/"
		} else {
			query = fmt.Sprintf("/?count=%d", *count)
		}
	} else {
		if count == nil {
			query = fmt.Sprintf("/?offset=%d", *offset)
		} else {
			query = fmt.Sprintf("/?offset=%d&count=%d", *offset, *count)
		}
	}

	request := httptest.NewRequest(http.MethodGet, query, io.NopCloser(bytes.NewBuffer([]byte{})))
	request.Header.Set("Authorization", TOKEN_STRING)
	return request
}

func TestFileDataEndpointShouldReturnProperHeadersDuringPreflight(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	defer databaseMock.AssertAllExpectionsSatisfied()
	jwtManagerMock := mock.NewJwtManagerMock(t)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()
	sut := endpoint.NewFileDataEndpoint(&databaseMock, &jwtManagerMock)

	router, responseRecorder := prepareGin()
	router.OPTIONS("/", sut.Options)
	request := httptest.NewRequest(http.MethodOptions, "/", io.NopCloser(bytes.NewReader([]byte{})))
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 200 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "" {
		fmt.Print("Expected body to be empty")
		t.FailNow()
	}
	if responseRecorder.Result().Header.Get("Access-Control-Allow-Headers") != "Authorization" {
		t.Error("Missing/Invalid 'Access-Control-Allow-Headers'")
	}
	if responseRecorder.Result().Header.Get("Access-Control-Allow-Methods") != "GET" {
		t.Error("Missing/Invalid 'Access-Control-Allow-Methods'")
	}
}

func TestFileDataEndpointShouldReturn500WhenQueryFailed(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectQuery(FILE_DATA_QUERY, fileDataToRequestArgs(FILE_DATA), []any{USER_ID, int64(15), int64(10)}, errors.New("db error"))
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{
		UserId: USER_ID,
	}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewFileDataEndpoint(&databaseMock, &jwtManagerMock)
	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)

	router.ServeHTTP(responseRecorder, prepareRequest(&OFFSET, &COUNT))

	if responseRecorder.Code != 500 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.Len() != 0 {
		t.Error("Expected empty response")
	}
}

func TestFileDataEndpointShouldReturn403WhenTokenIsInvalid(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID}, errors.New("token is invalid"))
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewFileDataEndpoint(&databaseMock, &jwtManagerMock)
	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)

	router.ServeHTTP(responseRecorder, prepareRequest(&OFFSET, &COUNT))

	if responseRecorder.Code != 403 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.Len() != 0 {
		t.Error("Expected empty response")
	}
}

func TestFileDataEndpointShouldReturn400WhenCountIsMissing(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewFileDataEndpoint(&databaseMock, &jwtManagerMock)
	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)

	router.ServeHTTP(responseRecorder, prepareRequest(&OFFSET, nil))

	if responseRecorder.Code != 400 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.Len() != 0 {
		t.Error("Expected empty response")
	}
}

func TestFileDataEndpointShouldReturn400WhenOffsetIsNegative(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewFileDataEndpoint(&databaseMock, &jwtManagerMock)
	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)

	router.ServeHTTP(responseRecorder, prepareRequest(&NEGATIVE_OFFSET, &COUNT))

	if responseRecorder.Code != 400 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.Len() != 0 {
		t.Error("Expected empty response")
	}
}

func TestFileDataEndpointShouldReturn400WhenCountIsNegative(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewFileDataEndpoint(&databaseMock, &jwtManagerMock)
	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)

	router.ServeHTTP(responseRecorder, prepareRequest(&OFFSET, &NEGATIVE_COUNT))

	if responseRecorder.Code != 400 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.Len() != 0 {
		t.Error("Expected empty response")
	}
}

func TestFileDataEndpointShouldReturn400WhenOffsetIsMissing(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewFileDataEndpoint(&databaseMock, &jwtManagerMock)
	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)

	router.ServeHTTP(responseRecorder, prepareRequest(nil, &COUNT))

	if responseRecorder.Code != 400 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.Len() != 0 {
		t.Error("Expected empty response")
	}
}

func TestFileDataEndpointShouldReturn200AndFileData(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectQuery(FILE_DATA_QUERY, fileDataToRequestArgs(FILE_DATA), []any{USER_ID, int64(15), int64(10)}, nil)
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{
		UserId: USER_ID,
	}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewFileDataEndpoint(&databaseMock, &jwtManagerMock)
	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)

	router.ServeHTTP(responseRecorder, prepareRequest(&OFFSET, &COUNT))

	if responseRecorder.Code != 200 {
		t.Error(responseRecorder.Code)
	}

	var fileData []endpoint.FileData
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &fileData)
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(fileData, FILE_DATA) {
		t.Error("Expected different file data")
	}
}
