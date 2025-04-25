// Package password handles hashing passwords.
package password

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

var logger *log.Logger = log.New(os.Stdout, "[PasswordFacade]: ", log.LstdFlags)

// IPasswordFacade interface provides methods for hashing passwords.
//
// HashPassword hashes given password and returns hash, if error occured returns error instead.
//
// MatchHashToPassword checks whether given hash match given raw password.
type IPasswordFacade interface {
	HashPassword(password string) (string, error)
	MatchHashToPassword(hash string, password string) bool
}

// PasswordFacade struct implements IPasswordFacade interface.
type PasswordFacade struct {
}

// HashPassword overrides IPasswordFacade.HashPassword.
// For passwords longer than 72 bytes error will be returned.
func (PasswordFacade) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		logger.Print(err.Error())
	}
	logger.Print("Hashed password")
	return string(hash), err
}

// HashPassword overrides IPasswordFacade.MatchHashToPassword.
func (PasswordFacade) MatchHashToPassword(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
