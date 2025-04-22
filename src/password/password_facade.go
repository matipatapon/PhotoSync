package password

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

var logger *log.Logger = log.New(os.Stdout, "[PasswordFacade]: ", log.LstdFlags)

type IPasswordFacade interface {
	HashPassword(password string) (string, error)
	CompareHashToPassword(hash string, password string) bool
}

type PasswordFacade struct {
}

func (PasswordFacade) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		logger.Print(err.Error())
	}
	logger.Print("Hashed password")
	return string(hash), err
}

func (PasswordFacade) CompareHashToPassword(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
