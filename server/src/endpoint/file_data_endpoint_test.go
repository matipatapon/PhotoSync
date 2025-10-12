package endpoint_test

import (
	"bytes"
	"encoding/base64"
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

var FILE_DATA_QUERY string = "SELECT id, filename, TO_CHAR(creation_date, 'YYYY.MM.DD HH24:MI:SS') AS date, mime_type, size, thumbnail FROM files WHERE user_id = $1 AND TO_CHAR(creation_date, 'YYYY.MM.DD') ILIKE $2 || '%' ORDER BY id DESC, creation_date DESC"
var OFFSET int64 = 10
var COUNT int64 = 15
var NEGATIVE_OFFSET int64 = -10
var NEGATIVE_COUNT int64 = -15

var FILE_DATA_WITH_THUMBNAIL []endpoint.FileData = []endpoint.FileData{
	{
		Id:           "13",
		Filename:     "B2.jpg",
		CreationDate: "2023.01.04 12:30:00",
		MIMEType:     "image/jpeg",
		Size:         "1024",
		Thumbnail:    "TU9SVEFERUxLQQ==",
	},
	{
		Id:           "17",
		Filename:     "9S.jpg",
		CreationDate: "2022.04.06 14:33:12",
		MIMEType:     "image/jpeg",
		Size:         "5000",
		Thumbnail:    "TFVORUNaS0E=",
	},
}

var FILE_DATA_WITHOUT_THUMBNAIL []endpoint.FileData = []endpoint.FileData{
	{
		Id:           "13",
		Filename:     "B2.jpg",
		CreationDate: "2023.01.04 12:30:00",
		MIMEType:     "image/jpeg",
		Size:         "1024",
		Thumbnail:    "",
	},
}

func fileDataToRequestArgs(fileData []endpoint.FileData) [][]any {
	result := [][]any{}
	for _, fd := range fileData {
		row := []any{}
		id, _ := strconv.ParseInt(fd.Id, 10, 64)
		row = append(row, id)
		row = append(row, fd.Filename)
		row = append(row, fd.CreationDate)
		row = append(row, int16(metadata.StringToMIMEType(fd.MIMEType)))
		size, _ := strconv.ParseInt(fd.Size, 10, 64)
		row = append(row, size)

		if fd.Thumbnail == "" {
			row = append(row, nil)
		} else {
			thumbnail, _ := base64.StdEncoding.DecodeString(fd.Thumbnail)
			row = append(row, thumbnail)
		}

		result = append(result, row)
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

func prepareRequest(date *string) *http.Request {
	var query string
	if date == nil {
		query = "/"
	} else {
		query = fmt.Sprintf("/?date=%s", *date)
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

func TestFileDataEndpointShouldReturn400WhenDateIsInvalid(t *testing.T) {
	invalidDates := []string{
		"yyyy.mm.dd",
		"%",
		"2025-12-24",
		"12.06.2025",
		"a12.06.2025b",
		"",
	}
	for _, invalidDate := range invalidDates {
		databaseMock := mock.NewDatabaseMock(t)
		defer databaseMock.AssertAllExpectionsSatisfied()

		jwtManagerMock := mock.NewJwtManagerMock(t)
		jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{
			UserId: USER_ID,
		}, nil)
		defer jwtManagerMock.AssertAllExpectionsSatisfied()

		sut := endpoint.NewFileDataEndpoint(&databaseMock, &jwtManagerMock)
		router, responseRecorder := prepareGin()
		router.GET("/", sut.Get)

		router.ServeHTTP(responseRecorder, prepareRequest(&invalidDate))

		if responseRecorder.Code != 400 {
			t.Error(responseRecorder.Code)
		}
		if responseRecorder.Body.Len() != 0 {
			t.Error("Expected empty response")
		}
	}
}

func TestFileDataEndpointShouldReturn500WhenQueryFailed(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectQuery(FILE_DATA_QUERY, fileDataToRequestArgs(FILE_DATA_WITH_THUMBNAIL), []any{USER_ID, PARAM_DATE}, errors.New("db error"))
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{
		UserId: USER_ID,
	}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewFileDataEndpoint(&databaseMock, &jwtManagerMock)
	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)

	router.ServeHTTP(responseRecorder, prepareRequest(&PARAM_DATE))

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

	router.ServeHTTP(responseRecorder, prepareRequest(&PARAM_DATE))

	if responseRecorder.Code != 403 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.Len() != 0 {
		t.Error("Expected empty response")
	}
}

func TestFileDataEndpointShouldReturn500WhenFileGotDeleted(t *testing.T) {
	results := [][][]any{{{}}, {}}
	for _, result := range results {
		fileId, _ := strconv.ParseInt(FILE_DATA_WITHOUT_THUMBNAIL[0].Id, 10, 64)

		databaseMock := mock.NewDatabaseMock(t)
		databaseMock.ExpectQuery(FILE_DATA_QUERY, fileDataToRequestArgs(FILE_DATA_WITHOUT_THUMBNAIL), []any{USER_ID, PARAM_DATE}, nil)
		databaseMock.ExpectQuery("SELECT file FROM files WHERE id = $1", result, []any{fileId}, nil)
		defer databaseMock.AssertAllExpectionsSatisfied()

		jwtManagerMock := mock.NewJwtManagerMock(t)
		jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{
			UserId: USER_ID,
		}, nil)
		defer jwtManagerMock.AssertAllExpectionsSatisfied()

		sut := endpoint.NewFileDataEndpoint(&databaseMock, &jwtManagerMock)
		router, responseRecorder := prepareGin()
		router.GET("/", sut.Get)

		router.ServeHTTP(responseRecorder, prepareRequest(&PARAM_DATE))

		if responseRecorder.Code != 500 {
			t.Error(responseRecorder.Code)
		}
		if responseRecorder.Body.Len() != 0 {
			t.Error("Expected empty response")
		}
	}
}

func TestFileDataEndpointShouldReturn500WhenFailedToGetAImage(t *testing.T) {
	fileId, _ := strconv.ParseInt(FILE_DATA_WITHOUT_THUMBNAIL[0].Id, 10, 64)

	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectQuery(FILE_DATA_QUERY, fileDataToRequestArgs(FILE_DATA_WITHOUT_THUMBNAIL), []any{USER_ID, PARAM_DATE}, nil)
	databaseMock.ExpectQuery("SELECT file FROM files WHERE id = $1", [][]any{{}}, []any{fileId}, errors.New("failed to get a file"))
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{
		UserId: USER_ID,
	}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewFileDataEndpoint(&databaseMock, &jwtManagerMock)
	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)

	router.ServeHTTP(responseRecorder, prepareRequest(&PARAM_DATE))

	if responseRecorder.Code != 500 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.Len() != 0 {
		t.Error("Expected empty response")
	}
}

func TestFileDataEndpointShouldReturnImageAsThumbnailWhenThereIsNoThumbnail(t *testing.T) {
	expectedFileData := make([]endpoint.FileData, len(FILE_DATA_WITHOUT_THUMBNAIL))
	copy(expectedFileData, FILE_DATA_WITHOUT_THUMBNAIL)
	expectedFileData[0].Thumbnail = base64.StdEncoding.EncodeToString(FILE)
	fileId, _ := strconv.ParseInt(expectedFileData[0].Id, 10, 64)

	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectQuery(FILE_DATA_QUERY, fileDataToRequestArgs(FILE_DATA_WITHOUT_THUMBNAIL), []any{USER_ID, PARAM_DATE}, nil)
	databaseMock.ExpectQuery("SELECT file FROM files WHERE id = $1", [][]any{{FILE}}, []any{fileId}, nil)
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{
		UserId: USER_ID,
	}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewFileDataEndpoint(&databaseMock, &jwtManagerMock)
	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)

	router.ServeHTTP(responseRecorder, prepareRequest(&PARAM_DATE))

	if responseRecorder.Code != 200 {
		t.Error(responseRecorder.Code)
	}

	var result []endpoint.FileData
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &result)
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(result, expectedFileData) {
		t.Error("Expected different file data")
	}
}

func TestFileDataEndpointShouldReturn200AndFileData(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectQuery(FILE_DATA_QUERY, fileDataToRequestArgs(FILE_DATA_WITH_THUMBNAIL), []any{USER_ID, PARAM_DATE}, nil)
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{
		UserId: USER_ID,
	}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewFileDataEndpoint(&databaseMock, &jwtManagerMock)
	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)

	router.ServeHTTP(responseRecorder, prepareRequest(&PARAM_DATE))

	if responseRecorder.Code != 200 {
		t.Error(responseRecorder.Code)
	}

	var fileData []endpoint.FileData
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &fileData)
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(fileData, FILE_DATA_WITH_THUMBNAIL) {
		t.Error("Expected different file data")
	}
}
