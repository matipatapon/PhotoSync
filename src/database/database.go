package database

type IDataBase interface {
	Query() ([][]any, error)
}
