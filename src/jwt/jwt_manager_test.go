package jwt_test

import (
	"fmt"
	"photosync/src/jwt"
	"testing"
	"time"
)

var USER_NAME string = "user321"

func TestJwtManagerShouldCreateAndParseToken(t *testing.T) {
	sut := jwt.NewJwtManager()

	expirationTime := time.Now().Unix() + int64(time.Second)*5
	tokenString, err := sut.Create(jwt.JwtPayload{Username: USER_NAME, ExpirationTime: expirationTime})
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
	if payload.ExpirationTime != expirationTime {
		fmt.Printf("expirationTime mismatch '%d' != '%d'", payload.ExpirationTime, expirationTime)
		t.FailNow()
	}
}

func TestJwtManagerShouldReturnErrorWhenTokenIsInvalid(t *testing.T) {
	sut := jwt.NewJwtManager()
	_, err := sut.Decode("invalid stringToken")
	if err == nil {
		t.FailNow()
	}
}

func TestEachJwtManagerShouldGenerateItsOwnKey(t *testing.T) {
	jm1 := jwt.NewJwtManager()
	jm2 := jwt.NewJwtManager()

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
	sut := jwt.NewJwtManager()
	tokenString, err := sut.Create(jwt.JwtPayload{Username: USER_NAME, ExpirationTime: time.Now().Unix()})
	if err != nil {
		t.FailNow()
	}

	time.Sleep(time.Second)
	_, err = sut.Decode(tokenString)
	if err == nil {
		fmt.Print("token shall be expired")
		t.FailNow()
	}
}
