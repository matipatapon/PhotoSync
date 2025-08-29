// Package database handles database operations
package database

// IDataBase interface provides method to perform queries on database
//
// Query is used to execute queries which return result e.g.
//  - Query(SELECT * FROM users)
//  - Query(COUNT(SELECT * FROM users WHERE name = $1), "Tomek")
//
// Execute is used to perform queries which doesn't return result e.g.
//  - Execute("INSERT INTO users VALUES($1, $2)", 1, "Karol")
//  - Execute("UPDATE users SET name = $1 WHERE id = $2", "Karol", 1)
// Mixing Query and Execute will result in undefined behaviour
//
// Pass argumments to query using args parameter to prevent SQL injection
type IDataBase interface {
	Query(sql string, args ...any) ([][]any, error)
	Execute(sql string, args ...any) error
	InitDb() error
	DropDb() error
}
