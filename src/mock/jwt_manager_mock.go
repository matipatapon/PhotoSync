package mock

import (
	"fmt"
	"photosync/src/jwt"
	"reflect"
	"testing"
)

type JwtManagerMock struct {
	expectedCreateData  []jwt.JwtPayload
	expectedCreateToken []string
	expectedCreateError []error
	expectedDecodeData  []jwt.JwtPayload
	expectedDecodeToken []string
	expectedDecodeError []error
	t                   *testing.T
}

func NewJwtManagerMock(t *testing.T) JwtManagerMock {
	return JwtManagerMock{[]jwt.JwtPayload{}, []string{}, []error{}, []jwt.JwtPayload{}, []string{}, []error{}, t}
}

func (jm *JwtManagerMock) Create(data jwt.JwtPayload) (string, error) {
	if len(jm.expectedCreateData) <= 0 {
		fmt.Print("Unexpected Create!")
		jm.t.FailNow()
	}
	expectedData := jm.expectedCreateData[len(jm.expectedCreateData)-1]
	jm.expectedCreateData = jm.expectedCreateData[:len(jm.expectedCreateData)-1]

	if !reflect.DeepEqual(expectedData, data) {
		fmt.Printf("Unexpected Payload! %v != %v", expectedData, data)
		jm.t.FailNow()
	}

	tokenString := jm.expectedCreateToken[len(jm.expectedCreateToken)-1]
	jm.expectedCreateToken = jm.expectedCreateToken[:len(jm.expectedCreateToken)-1]
	err := jm.expectedCreateError[len(jm.expectedCreateError)-1]
	jm.expectedCreateError = jm.expectedCreateError[:len(jm.expectedCreateError)-1]

	return tokenString, err
}

func (jm *JwtManagerMock) Decode(tokenString string) (jwt.JwtPayload, error) {
	if len(jm.expectedDecodeToken) <= 0 {
		fmt.Print("Unexpected Decode!")
		jm.t.FailNow()
	}
	expectedToken := jm.expectedDecodeToken[len(jm.expectedDecodeToken)-1]
	jm.expectedDecodeToken = jm.expectedDecodeToken[:len(jm.expectedDecodeToken)-1]

	if tokenString != expectedToken {
		fmt.Print("Unexpected Token! Expected '%s', Got '%s'", expectedToken, tokenString)
		jm.t.FailNow()
	}

	payload := jm.expectedDecodeData[len(jm.expectedDecodeData)-1]
	jm.expectedDecodeData = jm.expectedDecodeData[:len(jm.expectedDecodeData)-1]
	err := jm.expectedDecodeError[len(jm.expectedDecodeError)-1]
	jm.expectedDecodeError = jm.expectedDecodeError[:len(jm.expectedDecodeError)-1]

	return payload, err
}

func (jm *JwtManagerMock) ExpectCreate(data jwt.JwtPayload, tokenString string, err error) {
	jm.expectedCreateData = append(jm.expectedCreateData, data)
	jm.expectedCreateToken = append(jm.expectedCreateToken, tokenString)
	jm.expectedCreateError = append(jm.expectedCreateError, err)
}

func (jm *JwtManagerMock) ExpectDecode(tokenString string, data jwt.JwtPayload, err error) {
	jm.expectedDecodeData = append(jm.expectedDecodeData, data)
	jm.expectedDecodeToken = append(jm.expectedDecodeToken, tokenString)
	jm.expectedDecodeError = append(jm.expectedDecodeError, err)
}

func (jm *JwtManagerMock) AssertAllExpectionsSatisfied() {
	if len(jm.expectedDecodeToken) != 0 || len(jm.expectedCreateData) != 0 {
		fmt.Print("Not all expections satisfied!")
		jm.t.FailNow()
	}
}
