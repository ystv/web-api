package utils

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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
func NewDB(ctx context.Context, config DatabaseConfig) (*pgx.Conn, error) {
	dbURI := fmt.Sprintf("dbname=%s host=%s user=%s password=%s port=%s sslmode=%s application_name=web-api",
		config.Name, config.Host, config.Username, config.Password, config.Port, config.SSLMode) // Build connection string

	conn, err := pgxpool.Connect(ctx, dbURI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer conn.Close()

	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return conn, nil
}

// Transact wraps transactions
func Transact(ctx context.Context, db *pgx.Conn, txFunc func(pgx.Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback(ctx) // err is non-nil; don't change it
		} else {
			err = tx.Commit(ctx) // err is nil; if Commit returns error update err
		}
	}()
	err = txFunc(tx)
	return err
}
