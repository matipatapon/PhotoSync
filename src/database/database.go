package database

type IDataBase interface {
	Query(sql string, args ...any) ([][]any, error)
}
