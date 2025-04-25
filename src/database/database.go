// Package database handles database operations
package database

// IDataBase interface provides method to perform queries on database
//
// Query method attempts to perform query on database.
//
// Parameters:
//  - sql - query string e.g. "SELECT * FROM users"
//  - args - parameters which will be inserted into query string, this field shall be
//    used instead inserting them directly into query to prevent SQL INJECTION e.g.
//    Query("INSERT INTO users VALUES($1, $2)", 1, "Karol")
// Return values:
//  - When database returns response it is returned as two dimensional slice e.g.
//    Query("SELECT id, name FROM users") can return [[1, "Tomek"], [2, "Karol"]], nil
//  - When database returns no rows, empty silice shall be returned e.g.
//    Query("INSERT INTO users VALUES()$1, $2)", 3, "Andrzej") will return [], nil
//  - When database returns error it will be returned as second variable, slice will be nil e.g.
//    Query("INVALID SQL") will return nil, error
type IDataBase interface {
	Query(sql string, args ...any) ([][]any, error)
}
