package user

import (
	"errors"
	"photosync/src/mock"
	"testing"
)

var HASH string = "HASH"
var USERNAME string = "USERNAME"
var PASSWORD string = "PASSWORD"
var REGISTER_QUERY string = "INSERT INTO users VALUES($1, $2)"
var LOGIN_QUERY string = "SELECT password FROM users WHERE username = $1"
var ERROR error = errors.New("ERROR")

func TestUserFacadeShouldRegisterUser(t *testing.T) {
	dbMock := mock.NewDatabaseMock(t)
	dbMock.ExpectQuery(REGISTER_QUERY, [][]any{}, []any{USERNAME, HASH}, nil)
	passwordFacadeMock := mock.NewPasswordFacadeMock(t)
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
	dbMock := mock.NewDatabaseMock(t)
	passwordFacadeMock := mock.NewPasswordFacadeMock(t)
	passwordFacadeMock.ExpectHashPassword(PASSWORD, "", ERROR)
	sut := NewUserFacade(&dbMock, &passwordFacadeMock)

	err := sut.RegisterUser(USERNAME, PASSWORD)
	if err != ERROR {
		t.Fail()
	}

	passwordFacadeMock.AssertAllExpectionsSatisfied()
}

func TestUserFacadeShouldReturnErrorWhenQueryFailed(t *testing.T) {
	dbMock := mock.NewDatabaseMock(t)
	dbMock.ExpectQuery(REGISTER_QUERY, [][]any{}, []any{USERNAME, HASH}, ERROR)
	passwordFacadeMock := mock.NewPasswordFacadeMock(t)
	passwordFacadeMock.ExpectHashPassword(PASSWORD, HASH, nil)
	sut := NewUserFacade(&dbMock, &passwordFacadeMock)

	err := sut.RegisterUser(USERNAME, PASSWORD)
	if err != ERROR {
		t.Fail()
	}

	dbMock.AssertAllExpectionsSatisfied()
	passwordFacadeMock.AssertAllExpectionsSatisfied()
}

func TestUserFacadeShouldLoginUser(t *testing.T) {
	dbMock := mock.NewDatabaseMock(t)
	dbMock.ExpectQuery(LOGIN_QUERY, [][]any{{HASH}}, []any{USERNAME}, nil)
	passwordFacadeMock := mock.NewPasswordFacadeMock(t)
	passwordFacadeMock.ExpectMatchHashToPassword(HASH, PASSWORD, true)
	sut := NewUserFacade(&dbMock, &passwordFacadeMock)

	err := sut.LoginUser(USERNAME, PASSWORD)

	if err != nil {
		t.Fail()
	}
	dbMock.AssertAllExpectionsSatisfied()
	passwordFacadeMock.AssertAllExpectionsSatisfied()
}

func TestUserFacadeShouldReturnErrorWhenGivenUserDoesNotExist(t *testing.T) {
	dbMock := mock.NewDatabaseMock(t)
	dbMock.ExpectQuery(LOGIN_QUERY, [][]any{{}}, []any{USERNAME}, nil)
	passwordFacadeMock := mock.NewPasswordFacadeMock(t)
	sut := NewUserFacade(&dbMock, &passwordFacadeMock)

	err := sut.LoginUser(USERNAME, PASSWORD)

	if err == nil {
		t.Fail()
	}
	dbMock.AssertAllExpectionsSatisfied()
	passwordFacadeMock.AssertAllExpectionsSatisfied()
}
