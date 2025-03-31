package database

type IDataBase interface {
	QueryRow() ([]any, error)
}
