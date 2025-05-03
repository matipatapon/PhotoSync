package mock

import (
	"fmt"
	"testing"
)

type PasswordFacadeMock struct {
	t                             *testing.T
	HashPasswordExpectedPasswords []string
	HashPasswordHashes            []string
	HashPasswordErrors            []error
	MatchHashToPasswordHashes     []string
	MatchHashToPasswordPasswords  []string
	MatchHashToPasswordResults    []bool
}

func NewPasswordFacadeMock(t *testing.T) PasswordFacadeMock {
	return PasswordFacadeMock{t, []string{}, []string{}, []error{}, []string{}, []string{}, []bool{}}
}

func (pfm *PasswordFacadeMock) ExpectHashPassword(password string, hash string, err error) {
	pfm.HashPasswordErrors = append(pfm.HashPasswordErrors, err)
	pfm.HashPasswordExpectedPasswords = append(pfm.HashPasswordExpectedPasswords, password)
	pfm.HashPasswordHashes = append(pfm.HashPasswordHashes, hash)
}

func (pfm *PasswordFacadeMock) ExpectMatchHashToPassword(hash string, password string, result bool) {
	pfm.MatchHashToPasswordHashes = append(pfm.MatchHashToPasswordHashes, hash)
	pfm.MatchHashToPasswordPasswords = append(pfm.MatchHashToPasswordPasswords, password)
	pfm.MatchHashToPasswordResults = append(pfm.MatchHashToPasswordResults, result)
}

func (pfm *PasswordFacadeMock) MatchHashToPassword(hash string, password string) bool {
	if len(pfm.MatchHashToPasswordResults) == 0 {
		fmt.Println("Unexpected call!")
		pfm.t.FailNow()
	}

	expectedHash := pfm.MatchHashToPasswordHashes[len(pfm.MatchHashToPasswordHashes)-1]
	pfm.MatchHashToPasswordHashes = pfm.MatchHashToPasswordHashes[:len(pfm.MatchHashToPasswordHashes)-1]
	if expectedHash != hash {
		fmt.Println("Unexpected hash!")
		pfm.t.FailNow()
	}

	expectedPassword := pfm.MatchHashToPasswordPasswords[len(pfm.MatchHashToPasswordPasswords)-1]
	pfm.MatchHashToPasswordPasswords = pfm.MatchHashToPasswordPasswords[:len(pfm.MatchHashToPasswordPasswords)-1]
	if expectedPassword != password {
		fmt.Println("Unexpected password!")
		pfm.t.FailNow()
	}

	result := pfm.MatchHashToPasswordResults[len(pfm.MatchHashToPasswordResults)-1]
	pfm.MatchHashToPasswordResults = pfm.MatchHashToPasswordResults[:len(pfm.MatchHashToPasswordResults)-1]

	return result
}

func (pfm *PasswordFacadeMock) HashPassword(password string) (string, error) {
	if len(pfm.HashPasswordExpectedPasswords) == 0 {
		fmt.Println("Unexpected call!")
		pfm.t.FailNow()
	}

	expectedPassword := pfm.HashPasswordExpectedPasswords[len(pfm.HashPasswordExpectedPasswords)-1]
	pfm.HashPasswordExpectedPasswords = pfm.HashPasswordExpectedPasswords[:len(pfm.HashPasswordExpectedPasswords)-1]
	if password != expectedPassword {
		fmt.Println("Unexpected password!")
		pfm.t.FailNow()
	}

	hash := pfm.HashPasswordHashes[len(pfm.HashPasswordHashes)-1]
	pfm.HashPasswordHashes = pfm.HashPasswordHashes[:len(pfm.HashPasswordHashes)-1]

	err := pfm.HashPasswordErrors[len(pfm.HashPasswordErrors)-1]
	pfm.HashPasswordErrors = pfm.HashPasswordErrors[:len(pfm.HashPasswordErrors)-1]

	return hash, err
}

func (pfm *PasswordFacadeMock) AssertAllExpectionsSatisfied() {
	if len(pfm.HashPasswordExpectedPasswords) != 0 || len(pfm.MatchHashToPasswordResults) != 0 {
		fmt.Println("Not all expections satisfied!")
		pfm.t.FailNow()
	}
}
