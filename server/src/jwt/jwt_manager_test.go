package jwt_test

import (
	"fmt"
	"photosync/src/jwt"
	"photosync/src/mock"
	"testing"
)

var NOT_EXPIRED_TIME int64 = 96
var EXPIRATION_TIME int64 = 100
var EXPIRED_TIME int64 = 102

var USER_NAME string = "user321"
var USER_ID int64 = 1

var BIG_EXPIRATION_TIME int64 = 9223372036854775806
var BIG_USER_ID int64 = 9223372036854775805

func TestJwtManagerShouldHandleBigIntegers(t *testing.T) {
	thMock := mock.NewTimeHelperMock(t)
	thMock.ExpectTimeNow(BIG_EXPIRATION_TIME)
	defer thMock.AssertAllExpectionsSatisfied()

	sut := jwt.NewJwtManager(&thMock)

	tokenString, err := sut.Create(jwt.JwtPayload{UserId: BIG_USER_ID, Username: USER_NAME, ExpirationTime: BIG_EXPIRATION_TIME})
	if err != nil {
		t.Fail()
	}

	payload, err := sut.Decode(tokenString)
	if err != nil {
		fmt.Printf("unexpected error %s", err.Error())
		t.FailNow()
	}
	if payload.Username != USER_NAME {
		fmt.Printf("username mismatch '%s' != '%s'", payload.Username, USER_NAME)
		t.FailNow()
	}
	if payload.ExpirationTime != BIG_EXPIRATION_TIME {
		fmt.Printf("expirationTime mismatch '%d' != '%d'", payload.ExpirationTime, BIG_EXPIRATION_TIME)
		t.FailNow()
	}
	if payload.UserId != BIG_USER_ID {
		fmt.Printf("user_id mismatch '%d' != '%d'", payload.UserId, BIG_USER_ID)
		t.FailNow()
	}
}

func TestJwtManagerShouldCreateAndParseToken(t *testing.T) {
	thMock := mock.NewTimeHelperMock(t)
	thMock.ExpectTimeNow(NOT_EXPIRED_TIME)
	defer thMock.AssertAllExpectionsSatisfied()

	sut := jwt.NewJwtManager(&thMock)

	tokenString, err := sut.Create(jwt.JwtPayload{UserId: USER_ID, Username: USER_NAME, ExpirationTime: EXPIRATION_TIME})
	if err != nil {
		t.Fail()
	}

	payload, err := sut.Decode(tokenString)
	if err != nil {
		fmt.Printf("unexpected error %s", err.Error())
		t.FailNow()
	}
	if payload.Username != USER_NAME {
		fmt.Printf("username mismatch '%s' != '%s'", payload.Username, USER_NAME)
		t.FailNow()
	}
	if payload.ExpirationTime != EXPIRATION_TIME {
		fmt.Printf("expirationTime mismatch '%d' != '%d'", payload.ExpirationTime, EXPIRATION_TIME)
		t.FailNow()
	}
	if payload.UserId != USER_ID {
		fmt.Printf("user_id mismatch '%d' != '%d'", payload.UserId, USER_ID)
		t.FailNow()
	}
}

func TestJwtManagerShouldReturnErrorWhenTokenIsInvalid(t *testing.T) {
	thMock := mock.NewTimeHelperMock(t)
	defer thMock.AssertAllExpectionsSatisfied()

	sut := jwt.NewJwtManager(&thMock)
	_, err := sut.Decode("invalid stringToken")
	if err == nil {
		t.FailNow()
	}
}

func TestEachJwtManagerShouldGenerateItsOwnKey(t *testing.T) {
	thMock := mock.NewTimeHelperMock(t)
	defer thMock.AssertAllExpectionsSatisfied()

	jm1 := jwt.NewJwtManager(&thMock)
	jm2 := jwt.NewJwtManager(&thMock)

	tokenString, err := jm1.Create(jwt.JwtPayload{Username: USER_NAME})
	if err != nil {
		t.FailNow()
	}

	_, err = jm2.Decode(tokenString)
	if err == nil {
		t.FailNow()
	}
}

func TestJwtManagerShouldReturnErrorWhenTokenExpired(t *testing.T) {
	thMock := mock.NewTimeHelperMock(t)
	thMock.ExpectTimeNow(EXPIRED_TIME)
	defer thMock.AssertAllExpectionsSatisfied()

	sut := jwt.NewJwtManager(&thMock)
	tokenString, err := sut.Create(jwt.JwtPayload{Username: USER_NAME, ExpirationTime: EXPIRATION_TIME})
	if err != nil {
		t.FailNow()
	}

	_, err = sut.Decode(tokenString)
	if err == nil {
		fmt.Print("token shall be expired")
		t.FailNow()
	}
}
