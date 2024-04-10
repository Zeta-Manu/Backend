package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Zeta-Manu/Backend/internal/config"
)

type Database struct {
	Conn *sql.DB
}

// NewDatabase creates a new MySQL database connection.
func NewDatabase(dataSourceName string) (*Database, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Check if the database connection is alive
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{Conn: db}, nil
}

// InitializeDatabase initializes and returns a new database connection.
func InitializeDatabase(dbConfig config.DatabaseConfig) (*Database, error) {
	dbDataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)
	return NewDatabase(dbDataSourceName)
}

// Close closes the database connection.
func (db *Database) Close() error {
	if db.Conn != nil {
		return db.Conn.Close()
	}
	return nil
}
