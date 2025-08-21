package mock

import (
	"fmt"
	"photosync/src/helper"
	"testing"
)

type EnvGetterMock struct {
	expectedNames  helper.List[string]
	expectedValues helper.List[string]
	t              *testing.T
}

func NewEnvGetterMock(t *testing.T) EnvGetterMock {
	return EnvGetterMock{t: t}
}

func (egm *EnvGetterMock) ExpectGet(name string, value string) {
	egm.expectedNames.Append(name)
	egm.expectedValues.Append(value)
}

func (egm *EnvGetterMock) Get(name string) string {
	if egm.expectedNames.Length() == 0 {
		fmt.Println("Unexpected Get!")
		egm.t.FailNow()
	}

	expectedName := egm.expectedNames.PopFirst()
	if expectedName != name {
		fmt.Println("Unexpected Name!")
		egm.t.FailNow()
	}

	return egm.expectedValues.PopFirst()
}

func (egm *EnvGetterMock) AssertAllExpectionsSatisfied() {
	if egm.expectedNames.Length() != 0 {
		fmt.Println("Not all expects satisfied!")
		egm.t.FailNow()
	}
}
