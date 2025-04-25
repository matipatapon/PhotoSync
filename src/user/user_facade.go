// Package user handles user related operations.
package user

import (
	"log"
	"os"
	"photosync/src/database"
	"photosync/src/password"
)

var logger *log.Logger = log.New(os.Stdout, "[UserFacade]: ", log.LstdFlags)

// IUserFacade interface provides user related methods
type IUserFacade interface {
	RegisterUser(name string, password string) error
}

// UserFacade struct implements IUserFacade.
type UserFacade struct {
	db             database.IDataBase
	passwordFacade password.IPasswordFacade
}

// NewUserFacade creates UserFacade.
func NewUserFacade(db database.IDataBase, passwordFacade password.IPasswordFacade) UserFacade {
	return UserFacade{db, passwordFacade}
}

// RegisterUser overrides IUserFacade.RegisterUser.
// It inserts user with hashed password into database.
//
// Error will be returned when:
//   - Failed to hash password.
//   - Failed to insert user into database.
func (uf *UserFacade) RegisterUser(name string, password string) error {
	hash, err := uf.passwordFacade.HashPassword(password)
	if err != nil {
		logger.Print("Failed to hash password!")
		return err
	}

	_, err = uf.db.Query("INSERT INTO users VALUES($1, $2)", name, hash)
	if err != nil {
		logger.Print("Failed to add user into db!")
		return err
	}
	logger.Printf("Registered %s", name)
	return nil
}
