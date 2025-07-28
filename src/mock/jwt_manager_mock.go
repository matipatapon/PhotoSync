package mock

import (
	"fmt"
	"photosync/src/helper"
	"photosync/src/jwt"
	"reflect"
	"testing"
)

type JwtManagerMock struct {
	expectedCreateData  helper.List[jwt.JwtPayload]
	expectedCreateToken helper.List[string]
	expectedCreateError helper.List[error]
	expectedDecodeData  helper.List[jwt.JwtPayload]
	expectedDecodeToken helper.List[string]
	expectedDecodeError helper.List[error]
	t                   *testing.T
}

func NewJwtManagerMock(t *testing.T) JwtManagerMock {
	return JwtManagerMock{t: t}
}

func (jm *JwtManagerMock) Create(data jwt.JwtPayload) (string, error) {
	if jm.expectedCreateData.Length() <= 0 {
		fmt.Print("Unexpected Create!")
		jm.t.FailNow()
	}

	expectedData := jm.expectedCreateData.PopFirst()
	if !reflect.DeepEqual(expectedData, data) {
		fmt.Printf("Unexpected Payload! %v != %v", expectedData, data)
		jm.t.FailNow()
	}

	return jm.expectedCreateToken.PopFirst(), jm.expectedCreateError.PopFirst()
}

func (jm *JwtManagerMock) Decode(tokenString string) (jwt.JwtPayload, error) {
	if jm.expectedDecodeToken.Length() <= 0 {
		fmt.Print("Unexpected Decode!")
		jm.t.FailNow()
	}

	expectedToken := jm.expectedDecodeToken.PopFirst()
	if tokenString != expectedToken {
		fmt.Printf("Unexpected Token! Expected '%s', Got '%s'", expectedToken, tokenString)
		jm.t.FailNow()
	}

	return jm.expectedDecodeData.PopFirst(), jm.expectedDecodeError.PopFirst()
}

func (jm *JwtManagerMock) ExpectCreate(data jwt.JwtPayload, tokenString string, err error) {
	jm.expectedCreateData.Append(data)
	jm.expectedCreateToken.Append(tokenString)
	jm.expectedCreateError.Append(err)
}

func (jm *JwtManagerMock) ExpectDecode(tokenString string, data jwt.JwtPayload, err error) {
	jm.expectedDecodeData.Append(data)
	jm.expectedDecodeToken.Append(tokenString)
	jm.expectedDecodeError.Append(err)
}

func (jm *JwtManagerMock) AssertAllExpectionsSatisfied() {
	if jm.expectedDecodeToken.Length() != 0 || jm.expectedCreateData.Length() != 0 {
		fmt.Print("Not all expections satisfied!")
		jm.t.FailNow()
	}
}
