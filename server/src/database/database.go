package database

type IDataBase interface {
	Query(sql string, args ...any) ([][]any, error)
	Execute(sql string, args ...any) error
	InitDb() error
	DropDb() error
}
