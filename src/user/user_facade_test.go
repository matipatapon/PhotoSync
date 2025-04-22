package user

import (
	"errors"
	"photosync/src/database"
	"photosync/src/password"
	"testing"
)

var HASH string = "HASH"
var USERNAME string = "USERNAME"
var PASSWORD string = "PASSWORD"
var QUERY string = "INSERT INTO users VALUES($1, $2)"
var ERROR error = errors.New("ERROR")

func TestUserFacadeShouldRegisterUser(t *testing.T) {
	dbMock := database.NewDatabaseMock(t)
	dbMock.ExpectQuery(QUERY, [][]any{}, []any{USERNAME, HASH}, nil)
	passwordFacadeMock := password.NewPasswordFacadeMock(t)
	passwordFacadeMock.ExpectHashPassword(PASSWORD, HASH, nil)
	sut := NewUserFacade(&dbMock, &passwordFacadeMock)

	err := sut.RegisterUser(USERNAME, PASSWORD)
	if err != nil {
		t.Fail()
	}

	passwordFacadeMock.AssertAllExpectionsSatisfied()
	dbMock.AssertAllExpectionsSatisfied()
}

func TestUserFacadeShouldReturnErrorWhenFailedToHashPassword(t *testing.T) {
	dbMock := database.NewDatabaseMock(t)
	passwordFacadeMock := password.NewPasswordFacadeMock(t)
	passwordFacadeMock.ExpectHashPassword(PASSWORD, "", ERROR)
	sut := NewUserFacade(&dbMock, &passwordFacadeMock)

	err := sut.RegisterUser(USERNAME, PASSWORD)
	if err != ERROR {
		t.Fail()
	}

	passwordFacadeMock.AssertAllExpectionsSatisfied()
}

func TestUserFacadeShouldReturnErrorWhenQueryFailed(t *testing.T) {
	dbMock := database.NewDatabaseMock(t)
	dbMock.ExpectQuery(QUERY, [][]any{}, []any{USERNAME, HASH}, ERROR)
	passwordFacadeMock := password.NewPasswordFacadeMock(t)
	passwordFacadeMock.ExpectHashPassword(PASSWORD, HASH, nil)
	sut := NewUserFacade(&dbMock, &passwordFacadeMock)

	err := sut.RegisterUser(USERNAME, PASSWORD)
	if err != ERROR {
		t.Fail()
	}

	dbMock.AssertAllExpectionsSatisfied()
	passwordFacadeMock.AssertAllExpectionsSatisfied()
}
