package user

import (
	"log"
	"os"
	"photosync/src/database"
	"photosync/src/password"
)

var logger *log.Logger = log.New(os.Stdout, "[UserFacade]: ", log.LstdFlags)

type UserFacade struct {
	db             database.IDataBase
	passwordFacade password.IPasswordFacade
}

func NewUserFacade(db database.IDataBase, passwordFacade password.IPasswordFacade) UserFacade {
	return UserFacade{db, passwordFacade}
}

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
