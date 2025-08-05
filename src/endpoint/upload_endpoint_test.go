package endpoint_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"photosync/src/endpoint"
	"photosync/src/jwt"
	"photosync/src/metadata"
	"photosync/src/mock"
	"strconv"
	"testing"
)

var UPLOAD_SQL string = "INSERT INTO files(user_id, creation_date, filename, mime_type, file, hash, size) VALUES($1, TO_TIMESTAMP($2, 'YYYY.MM.DD HH24:MI:SS'), $3, $4, $5, $6, $7) RETURNING id"

func createRequest(fields map[string][]byte, token string) *http.Request {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	for key, value := range fields {
		fw, _ := w.CreateFormField(key)
		io.Copy(fw, bytes.NewBuffer(value))
	}
	w.Close()

	request := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(body.Bytes())))
	request.Header.Set("Content-Type", w.FormDataContentType())
	request.Header.Set("Authorization", token)

	return request
}

func TestUploadEndpointShouldReturn202WhenImageAlreadyExistsInDb(t *testing.T) {
	queryResults := [][][]any{
		{},
		{{}},
	}
	for _, queryResult := range queryResults {
		databaseMock := mock.NewDatabaseMock(t)
		databaseMock.ExpectQuery(UPLOAD_SQL, queryResult, []any{USER_ID, MODIFICATION_DATE, FILENAME, metadata.JPG, FILE, HASH, len(FILE)}, nil)
		defer databaseMock.AssertAllExpectionsSatisfied()

		metadataExtractorMock := mock.NewMetadataExtractorMock(t)
		metadataExtractorMock.ExpectExtract(FILE, metadata.Metadata{MIMEType: metadata.JPG})
		defer metadataExtractorMock.AssertAllExpectionsSatisfied()

		hasherMock := mock.NewHasherMock(t)
		hasherMock.ExpectHash(FILE, HASH, nil)
		defer hasherMock.AssertAllExpectionsSatisfied()

		jwtManagerMock := mock.NewJwtManagerMock(t)
		jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: EXPIRATION_TIME}, nil)
		defer jwtManagerMock.AssertAllExpectionsSatisfied()

		sut := endpoint.NewUploadEndpoint(&databaseMock, &metadataExtractorMock, &hasherMock, &jwtManagerMock)

		request := createRequest(
			map[string][]byte{
				"filename":          []byte(FILENAME),
				"modification_date": []byte(MODIFICATION_DATE),
				"file":              FILE,
			},
			TOKEN_STRING,
		)

		router, responseRecorder := prepareGin()
		router.POST("/", sut.Post)
		router.ServeHTTP(responseRecorder, request)

		if responseRecorder.Code != 402 {
			t.Error(responseRecorder.Code)
		}
		if responseRecorder.Body.String() != "" {
			fmt.Print("Expected body to be empty")
			t.FailNow()
		}
	}
}

func TestUploadEndpointShouldReturn401ForUnsupportedFile(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	defer databaseMock.AssertAllExpectionsSatisfied()

	metadataExtractorMock := mock.NewMetadataExtractorMock(t)
	metadataExtractorMock.ExpectExtract(FILE, metadata.Metadata{MIMEType: metadata.UNKNOWN})
	defer metadataExtractorMock.AssertAllExpectionsSatisfied()

	hasherMock := mock.NewHasherMock(t)
	hasherMock.ExpectHash(FILE, HASH, nil)
	defer hasherMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: EXPIRATION_TIME}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewUploadEndpoint(&databaseMock, &metadataExtractorMock, &hasherMock, &jwtManagerMock)

	request := createRequest(
		map[string][]byte{
			"filename":          []byte(FILENAME),
			"modification_date": []byte(MODIFICATION_DATE),
			"file":              FILE,
		},
		TOKEN_STRING,
	)

	router, responseRecorder := prepareGin()
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 401 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "" {
		fmt.Print("Expected body to be empty")
		t.FailNow()
	}
}

func TestUploadEndpointShouldReturn400WhenModificationDateIsInvalid(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	defer databaseMock.AssertAllExpectionsSatisfied()

	metadataExtractorMock := mock.NewMetadataExtractorMock(t)
	defer metadataExtractorMock.AssertAllExpectionsSatisfied()

	hasherMock := mock.NewHasherMock(t)
	defer hasherMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: EXPIRATION_TIME}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewUploadEndpoint(&databaseMock, &metadataExtractorMock, &hasherMock, &jwtManagerMock)

	request := createRequest(
		map[string][]byte{
			"filename":          []byte(FILENAME),
			"modification_date": []byte(INVALID_MODIFICATION_DATE),
			"file":              FILE,
		},
		TOKEN_STRING,
	)

	router, responseRecorder := prepareGin()
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 400 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "" {
		fmt.Print("Expected body to be empty")
		t.FailNow()
	}
}

func TestUploadEndpointShouldPrioritizeCreationDateFromMetadata(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectQuery(UPLOAD_SQL, [][]any{{FILE_ID}}, []any{USER_ID, CREATION_DATE, FILENAME, metadata.JPG, FILE, HASH, len(FILE)}, nil)
	defer databaseMock.AssertAllExpectionsSatisfied()

	creationDate, _ := metadata.NewDate(CREATION_DATE)
	metadataExtractorMock := mock.NewMetadataExtractorMock(t)
	metadataExtractorMock.ExpectExtract(FILE, metadata.Metadata{MIMEType: metadata.JPG, CreationDate: &creationDate})
	defer metadataExtractorMock.AssertAllExpectionsSatisfied()

	hasherMock := mock.NewHasherMock(t)
	hasherMock.ExpectHash(FILE, HASH, nil)
	defer hasherMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: EXPIRATION_TIME}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewUploadEndpoint(&databaseMock, &metadataExtractorMock, &hasherMock, &jwtManagerMock)

	request := createRequest(
		map[string][]byte{
			"filename":          []byte(FILENAME),
			"modification_date": []byte(MODIFICATION_DATE),
			"file":              FILE,
		},
		TOKEN_STRING,
	)

	router, responseRecorder := prepareGin()
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 200 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.String() != strconv.FormatInt(FILE_ID, 10) {
		fmt.Printf("Expected '%s', got '%s'", strconv.FormatInt(FILE_ID, 10), responseRecorder.Body.String())
		t.FailNow()
	}
}

func TestUploadEndpointShouldReturn500WhenFailedToSaveFileInDatabase(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectQuery(UPLOAD_SQL, [][]any{}, []any{USER_ID, MODIFICATION_DATE, FILENAME, metadata.JPG, FILE, HASH, len(FILE)}, errors.New("DB error"))
	defer databaseMock.AssertAllExpectionsSatisfied()

	metadataExtractorMock := mock.NewMetadataExtractorMock(t)
	metadataExtractorMock.ExpectExtract(FILE, metadata.Metadata{MIMEType: metadata.JPG})
	defer metadataExtractorMock.AssertAllExpectionsSatisfied()

	hasherMock := mock.NewHasherMock(t)
	hasherMock.ExpectHash(FILE, HASH, nil)
	defer hasherMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: EXPIRATION_TIME}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewUploadEndpoint(&databaseMock, &metadataExtractorMock, &hasherMock, &jwtManagerMock)

	request := createRequest(
		map[string][]byte{
			"filename":          []byte(FILENAME),
			"modification_date": []byte(MODIFICATION_DATE),
			"file":              FILE,
		},
		TOKEN_STRING,
	)

	router, responseRecorder := prepareGin()
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 500 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "" {
		fmt.Print("Expected body to be empty")
		t.FailNow()
	}
}

func TestUploadEndpointShouldReturn500WhenFailedToHashAFile(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	defer databaseMock.AssertAllExpectionsSatisfied()

	metadataExtractorMock := mock.NewMetadataExtractorMock(t)
	defer metadataExtractorMock.AssertAllExpectionsSatisfied()

	hasherMock := mock.NewHasherMock(t)
	hasherMock.ExpectHash(FILE, "", errors.New("Failed to hash a file"))
	defer hasherMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: EXPIRATION_TIME}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewUploadEndpoint(&databaseMock, &metadataExtractorMock, &hasherMock, &jwtManagerMock)

	request := createRequest(
		map[string][]byte{
			"filename":          []byte(FILENAME),
			"modification_date": []byte(MODIFICATION_DATE),
			"file":              FILE,
		},
		TOKEN_STRING,
	)

	router, responseRecorder := prepareGin()
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 500 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "" {
		fmt.Print("Expected body to be empty")
		t.FailNow()
	}
}

func TestUploadEndpointShouldReturn400WhenAnyRequestPartIsMissing(t *testing.T) {
	requests := []*http.Request{
		createRequest(
			map[string][]byte{
				"modification_date": []byte(MODIFICATION_DATE),
				"file":              FILE,
			}, TOKEN_STRING),
		createRequest(
			map[string][]byte{
				"filename": []byte(FILENAME),
				"file":     FILE,
			}, TOKEN_STRING),
		createRequest(
			map[string][]byte{
				"filename":          []byte(FILENAME),
				"modification_date": []byte(MODIFICATION_DATE),
			}, TOKEN_STRING),
	}

	for _, request := range requests {
		databaseMock := mock.NewDatabaseMock(t)
		defer databaseMock.AssertAllExpectionsSatisfied()

		metadataExtractorMock := mock.NewMetadataExtractorMock(t)
		defer metadataExtractorMock.AssertAllExpectionsSatisfied()

		hasherMock := mock.NewHasherMock(t)
		defer hasherMock.AssertAllExpectionsSatisfied()

		jwtManagerMock := mock.NewJwtManagerMock(t)
		jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: EXPIRATION_TIME}, nil)
		defer jwtManagerMock.AssertAllExpectionsSatisfied()

		sut := endpoint.NewUploadEndpoint(&databaseMock, &metadataExtractorMock, &hasherMock, &jwtManagerMock)

		router, responseRecorder := prepareGin()
		router.POST("/", sut.Post)
		router.ServeHTTP(responseRecorder, request)

		if responseRecorder.Code != 400 {
			t.Error(responseRecorder.Code)
		}
		if responseRecorder.Body.String() != "" {
			fmt.Print("Expected body to be empty")
			t.FailNow()
		}
	}
}

func TestUploadEndpointShouldReturn403WhenInvalidJwt(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	defer databaseMock.AssertAllExpectionsSatisfied()

	metadataExtractorMock := mock.NewMetadataExtractorMock(t)
	defer metadataExtractorMock.AssertAllExpectionsSatisfied()

	hasherMock := mock.NewHasherMock(t)
	defer hasherMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{}, errors.New("JWT is invalid"))
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewUploadEndpoint(&databaseMock, &metadataExtractorMock, &hasherMock, &jwtManagerMock)

	request := createRequest(
		map[string][]byte{
			"filename":          []byte(FILENAME),
			"modification_date": []byte(MODIFICATION_DATE),
			"file":              FILE,
		},
		TOKEN_STRING,
	)

	router, responseRecorder := prepareGin()
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 403 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "" {
		fmt.Print("Expected body to be empty")
		t.FailNow()
	}
}

func TestUploadEndpointShouldReturn400ForEmptyRequestBody(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	defer databaseMock.AssertAllExpectionsSatisfied()
	metadataExtractorMock := mock.NewMetadataExtractorMock(t)
	defer metadataExtractorMock.AssertAllExpectionsSatisfied()
	hasherMock := mock.NewHasherMock(t)
	defer hasherMock.AssertAllExpectionsSatisfied()
	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: EXPIRATION_TIME}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()
	sut := endpoint.NewUploadEndpoint(&databaseMock, &metadataExtractorMock, &hasherMock, &jwtManagerMock)

	request := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader([]byte{})))
	request.Header.Set("Authorization", TOKEN_STRING)

	router, responseRecorder := prepareGin()
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 400 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.String() != "" {
		fmt.Print("Expected body to be empty")
		t.FailNow()
	}
}

func TestUploadEndpointShouldSaveGivenImageToDb(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectQuery(UPLOAD_SQL, [][]any{{FILE_ID}}, []any{USER_ID, MODIFICATION_DATE, FILENAME, metadata.JPG, FILE, HASH, len(FILE)}, nil)
	defer databaseMock.AssertAllExpectionsSatisfied()

	metadataExtractorMock := mock.NewMetadataExtractorMock(t)
	metadataExtractorMock.ExpectExtract(FILE, metadata.Metadata{MIMEType: metadata.JPG})
	defer metadataExtractorMock.AssertAllExpectionsSatisfied()

	hasherMock := mock.NewHasherMock(t)
	hasherMock.ExpectHash(FILE, HASH, nil)
	defer hasherMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: EXPIRATION_TIME}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewUploadEndpoint(&databaseMock, &metadataExtractorMock, &hasherMock, &jwtManagerMock)

	request := createRequest(
		map[string][]byte{
			"filename":          []byte(FILENAME),
			"modification_date": []byte(MODIFICATION_DATE),
			"file":              FILE,
		},
		TOKEN_STRING,
	)

	router, responseRecorder := prepareGin()
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 200 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.String() != strconv.FormatInt(FILE_ID, 10) {
		fmt.Printf("Expected '%s', got '%s'", strconv.FormatInt(FILE_ID, 10), responseRecorder.Body.String())
		t.FailNow()
	}
}

func TestUploadEndpointShouldHandleRequestPartsInDifferentOrder(t *testing.T) {
	databaseMock := mock.NewDatabaseMock(t)
	databaseMock.ExpectQuery(UPLOAD_SQL, [][]any{{FILE_ID}}, []any{USER_ID, MODIFICATION_DATE, FILENAME, metadata.JPG, FILE, HASH, len(FILE)}, nil)
	defer databaseMock.AssertAllExpectionsSatisfied()

	metadataExtractorMock := mock.NewMetadataExtractorMock(t)
	metadataExtractorMock.ExpectExtract(FILE, metadata.Metadata{MIMEType: metadata.JPG})
	defer metadataExtractorMock.AssertAllExpectionsSatisfied()

	hasherMock := mock.NewHasherMock(t)
	hasherMock.ExpectHash(FILE, HASH, nil)
	defer hasherMock.AssertAllExpectionsSatisfied()

	jwtManagerMock := mock.NewJwtManagerMock(t)
	jwtManagerMock.ExpectDecode(TOKEN_STRING, jwt.JwtPayload{UserId: USER_ID, Username: USERNAME, ExpirationTime: EXPIRATION_TIME}, nil)
	defer jwtManagerMock.AssertAllExpectionsSatisfied()

	sut := endpoint.NewUploadEndpoint(&databaseMock, &metadataExtractorMock, &hasherMock, &jwtManagerMock)

	request := createRequest(
		map[string][]byte{
			"file":              FILE,
			"filename":          []byte(FILENAME),
			"modification_date": []byte(MODIFICATION_DATE),
		},
		TOKEN_STRING,
	)

	router, responseRecorder := prepareGin()
	router.POST("/", sut.Post)
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != 200 {
		t.Error(responseRecorder.Code)
	}
	if responseRecorder.Body.String() != strconv.FormatInt(FILE_ID, 10) {
		fmt.Printf("Expected '%s', got '%s'", strconv.FormatInt(FILE_ID, 10), responseRecorder.Body.String())
		t.FailNow()
	}
}
