package mock

import (
	"fmt"
	"photosync/src/helper"
	"reflect"
	"testing"
)

type HasherMock struct {
	expectedFiles  helper.List[[]byte]
	expectedHashes helper.List[string]
	expectedErrors helper.List[error]
	t              *testing.T
}

func NewHasherMock(t *testing.T) HasherMock {
	return HasherMock{t: t}
}

func (hm *HasherMock) ExpectHash(file []byte, hash string, err error) {
	hm.expectedFiles.Append(file)
	hm.expectedHashes.Append(hash)
	hm.expectedErrors.Append(err)
}

func (hm *HasherMock) Hash(file []byte) (string, error) {
	if hm.expectedFiles.Length() == 0 {
		fmt.Println("Unexpected Hash()")
		hm.t.FailNow()
	}
	if !reflect.DeepEqual(hm.expectedFiles.PopFirst(), file) {
		fmt.Println("Unexpected file")
		hm.t.FailNow()
	}

	return hm.expectedHashes.PopFirst(), hm.expectedErrors.PopFirst()
}

func (hm *HasherMock) AssertAllExpectionsSatisfied() {
	if hm.expectedFiles.Length() != 0 {
		fmt.Print("Not all expections satisfied!")
		hm.t.FailNow()
	}
}
