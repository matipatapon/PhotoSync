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
	"photosync/src/mock"
	"reflect"
	"strconv"
	"testing"
)

var DATES_QUERY string = "SELECT TO_CHAR(creation_date, 'YYYY.MM.DD') AS date, COUNT(*) AS file_count FROM files WHERE user_id = $1 GROUP BY date ORDER BY date DESC"
var DATES_QUERY_WITH_YEAR_FILTRATION string = "SELECT TO_CHAR(creation_date, 'YYYY.MM.DD') AS date, COUNT(*) AS file_count FROM files WHERE user_id = $1 AND DATE_PART('year', creation_date) = $2 GROUP BY date ORDER BY date DESC"
var DATES_QUERY_WITH_YEAR_AND_MONTH_FILTRATION string = "SELECT TO_CHAR(creation_date, 'YYYY.MM.DD') AS date, COUNT(*) AS file_count FROM files WHERE user_id = $1 AND DATE_PART('year', creation_date) = $2 AND DATE_PART('month', creation_date) = $3 GROUP BY date ORDER BY date DESC"

func datesToQueryResult(dates []endpoint.Date) [][]any {
	var result [][]any = [][]any{}
	for _, date := range dates {
		fileCount, _ := strconv.ParseInt(date.FileCount, 10, 64)
		result = append(result, []any{date.Date, fileCount})
	}
	return result
}

func TestDatesEndpointShouldReturnAllDates(t *testing.T) {
	dates := []endpoint.Date{{Date: "2025.05.16", FileCount: "12"}, {Date: "2024.06.13", FileCount: "100"}}
	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectQuery(DATES_QUERY, datesToQueryResult(dates), []any{USER_ID}, nil)
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: EXPIRATION_TIME}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewDatesEndpoint(&databaseMock, &jwtManagerMock)

	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)
	request := httptest.NewRequest(http.MethodGet, "/", io.NopCloser(bytes.NewReader([]byte{})))
	request.Header.Set("Authorization", TOKEN_STRING)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 200 {
		t.Error(responseRecorder.Code)
	}

	var resultDates []endpoint.Date
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &resultDates)
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(resultDates, dates) {
		t.Error("Expected different dates")
	}
}

func TestDatesEndpointShouldReturnEmptyArrayWhenThereAreNoDates(t *testing.T) {
	dates := []endpoint.Date{}
	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectQuery(DATES_QUERY, datesToQueryResult(dates), []any{USER_ID}, nil)
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: EXPIRATION_TIME}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewDatesEndpoint(&databaseMock, &jwtManagerMock)

	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)
	request := httptest.NewRequest(http.MethodGet, "/", io.NopCloser(bytes.NewReader([]byte{})))
	request.Header.Set("Authorization", TOKEN_STRING)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 200 {
		t.Error(responseRecorder.Code)
	}

	var resultDates []endpoint.Date
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &resultDates)
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(resultDates, dates) {
		t.Error("Expected different dates")
	}
}

func TestDatesEndpointShouldReturn403WhenTokenIsInvalid(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: EXPIRATION_TIME}, errors.New("invalid token"))
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewDatesEndpoint(&databaseMock, &jwtManagerMock)

	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)
	request := httptest.NewRequest(http.MethodGet, "/", io.NopCloser(bytes.NewReader([]byte{})))
	request.Header.Set("Authorization", TOKEN_STRING)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 403 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.Len() != 0 {
		t.Error("Expected empty body")
	}
}

func TestDatesEndpointShouldReturn500WhenQueryFailed(t *testing.T) {
	dates := []endpoint.Date{}
	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectQuery(DATES_QUERY, datesToQueryResult(dates), []any{USER_ID}, errors.New("query failed"))
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: EXPIRATION_TIME}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewDatesEndpoint(&databaseMock, &jwtManagerMock)

	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)
	request := httptest.NewRequest(http.MethodGet, "/", io.NopCloser(bytes.NewReader([]byte{})))
	request.Header.Set("Authorization", TOKEN_STRING)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 500 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.Len() != 0 {
		t.Error("Expected empty body")
	}
}

func TestDatesEndpointShouldFilterDatesByYear(t *testing.T) {
	year := int64(2025)
	dates := []endpoint.Date{}

	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectQuery(DATES_QUERY_WITH_YEAR_FILTRATION, datesToQueryResult(dates), []any{USER_ID, year}, nil)
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: EXPIRATION_TIME}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewDatesEndpoint(&databaseMock, &jwtManagerMock)

	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?year=%d", year), io.NopCloser(bytes.NewReader([]byte{})))
	request.Header.Set("Authorization", TOKEN_STRING)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 200 {
		t.Error(responseRecorder.Code)
	}

	var resultDates []endpoint.Date
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &resultDates)
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(resultDates, dates) {
		t.Error("Expected different dates")
	}
}

func TestDatesEndpointShouldReturn400WhenYearIsInvalid(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: EXPIRATION_TIME}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewDatesEndpoint(&databaseMock, &jwtManagerMock)

	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)
	request := httptest.NewRequest(http.MethodGet, "/?year=something", io.NopCloser(bytes.NewReader([]byte{})))
	request.Header.Set("Authorization", TOKEN_STRING)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 400 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.Len() != 0 {
		t.Error("Expected empty body")
	}
}

func TestDatesEndpointShouldFilterDatesByYearAndMonth(t *testing.T) {
	year := int64(2025)
	month := int64(6)
	dates := []endpoint.Date{}

	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectQuery(DATES_QUERY_WITH_YEAR_AND_MONTH_FILTRATION, datesToQueryResult(dates), []any{USER_ID, year, month}, nil)
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: EXPIRATION_TIME}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewDatesEndpoint(&databaseMock, &jwtManagerMock)

	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?year=%d&month=%d", year, month), io.NopCloser(bytes.NewReader([]byte{})))
	request.Header.Set("Authorization", TOKEN_STRING)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 200 {
		t.Error(responseRecorder.Code)
	}

	var resultDates []endpoint.Date
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &resultDates)
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(resultDates, dates) {
		t.Error("Expected different dates")
	}
}

func TestDatesEndpointShouldReturn400WhenMonthIsInvalid(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: EXPIRATION_TIME}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewDatesEndpoint(&databaseMock, &jwtManagerMock)

	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)
	request := httptest.NewRequest(http.MethodGet, "/?year=2025&month=something", io.NopCloser(bytes.NewReader([]byte{})))
	request.Header.Set("Authorization", TOKEN_STRING)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 400 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.Len() != 0 {
		t.Error("Expected empty body")
	}
}

func TestDatesEndpointShouldReturn400WhenMonthIsSpecifiedWithoutYear(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	defer databaseMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: EXPIRATION_TIME}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewDatesEndpoint(&databaseMock, &jwtManagerMock)

	router, responseRecorder := prepareGin()
	router.GET("/", sut.Get)
	request := httptest.NewRequest(http.MethodGet, "/?month=6", io.NopCloser(bytes.NewReader([]byte{})))
	request.Header.Set("Authorization", TOKEN_STRING)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 400 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.Len() != 0 {
		t.Error("Expected empty body")
	}
}
