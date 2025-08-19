package mock

import (
	"fmt"
	"photosync/src/helper"
	"testing"
)

type PasswordFacadeMock struct {
	t                             *testing.T
	HashPasswordExpectedPasswords helper.List[string]
	HashPasswordHashes            helper.List[string]
	HashPasswordErrors            helper.List[error]
	MatchHashToPasswordHashes     helper.List[string]
	MatchHashToPasswordPasswords  helper.List[string]
	MatchHashToPasswordResults    helper.List[bool]
}

func NewPasswordFacadeMock(t *testing.T) PasswordFacadeMock {
	return PasswordFacadeMock{t: t}
}

func (pfm *PasswordFacadeMock) ExpectHashPassword(password string, hash string, err error) {
	pfm.HashPasswordErrors.Append(err)
	pfm.HashPasswordExpectedPasswords.Append(password)
	pfm.HashPasswordHashes.Append(hash)
}

func (pfm *PasswordFacadeMock) ExpectMatchHashToPassword(hash string, password string, result bool) {
	pfm.MatchHashToPasswordHashes.Append(hash)
	pfm.MatchHashToPasswordPasswords.Append(password)
	pfm.MatchHashToPasswordResults.Append(result)
}

func (pfm *PasswordFacadeMock) MatchHashToPassword(hash string, password string) bool {
	if pfm.MatchHashToPasswordResults.Length() == 0 {
		fmt.Println("Unexpected call!")
		pfm.t.FailNow()
	}

	expectedHash := pfm.MatchHashToPasswordHashes.PopFirst()
	if expectedHash != hash {
		fmt.Println("Unexpected hash!")
		pfm.t.FailNow()
	}

	expectedPassword := pfm.MatchHashToPasswordPasswords.PopFirst()
	if expectedPassword != password {
		fmt.Println("Unexpected password!")
		pfm.t.FailNow()
	}

	return pfm.MatchHashToPasswordResults.PopFirst()
}

func (pfm *PasswordFacadeMock) HashPassword(password string) (string, error) {
	if pfm.HashPasswordExpectedPasswords.Length() == 0 {
		fmt.Println("Unexpected call!")
		pfm.t.FailNow()
	}

	expectedPassword := pfm.HashPasswordExpectedPasswords.PopFirst()
	if password != expectedPassword {
		fmt.Println("Unexpected password!")
		pfm.t.FailNow()
	}

	return pfm.HashPasswordHashes.PopFirst(), pfm.HashPasswordErrors.PopFirst()
}

func (pfm *PasswordFacadeMock) AssertAllExpectionsSatisfied() {
	if pfm.HashPasswordExpectedPasswords.Length() != 0 || pfm.MatchHashToPasswordResults.Length() != 0 {
		fmt.Println("Not all expections satisfied!")
		pfm.t.FailNow()
	}
}
