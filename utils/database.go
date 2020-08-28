package utils

import (
	"fmt"
	"log"
	"os"

	// PostgreSQL driver
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

// InitDB Initialises the connection to the database
func InitDB() *sqlx.DB {
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")

	dbURI := fmt.Sprintf("dbname=%s host=%s user=%s password=%s port=%s sslmode=disable", dbName, dbHost, username, password, dbPort) // Build connection string

	// Declared err since DB would be nil reference for when it is used outside, the := needed to be = essentially
	var err error
	db, err := sqlx.Open("postgres", dbURI)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Printf("Connected to DB: %s@%s", dbName, dbHost)
	return db
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
			tx.Rollback() // err is non-nil; don't chang eit
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()
	err = txFunc(tx)
	return err
}
