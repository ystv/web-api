package utils

import (
	"fmt"

	// PostgreSQL driver
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

// DatabaseConfig represents a configuration to connect to an SQL database
type DatabaseConfig struct {
	Host     string
	Port     string
	SSLMode  string
	Name     string
	Username string
	Password string
}

// NewDB Initialises the connection to the database
func NewDB(config DatabaseConfig) (*sqlx.DB, error) {
	dbURI := fmt.Sprintf("dbname=%s host=%s user=%s password=%s port=%s sslmode=%s application_name=web-api",
		config.Name, config.Host, config.Username, config.Password, config.Port, config.SSLMode) // Build connection string

	db, err := sqlx.Open("postgres", dbURI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}

// Transact wraps transactions
func Transact(db *sqlx.DB, txFunc func(*sqlx.Tx) error) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()
	err = txFunc(tx)
	return err
}
