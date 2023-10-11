package database

import (
	"database/sql"
)

type DBAdapter interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// Executes a SQL query and returns a result
func (db *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Execute a SQL statement
func (db *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := db.Conn.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}
