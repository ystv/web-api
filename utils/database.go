package utils

import (
	"fmt"
	"log"
	"os"

	// PostgreSQL driver
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

// DB object
var DB *sqlx.DB

// InitDB Initialises the connection to the database
func InitDB() {
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")

	dbURI := fmt.Sprintf("dbname=%s host=%s user=%s password=%s port=%s sslmode=disable", dbName, dbHost, username, password, dbPort) // Build connection string

	// Declared err since DB would be nil reference for when it is used outside, the := needed to be = essentially
	var err error
	DB, err = sqlx.Open("postgres", dbURI)
	if err != nil {
		panic(err)
	}
	err = DB.Ping()
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}

	log.Printf("Connected to DB: %s@%s", dbName, dbHost)
}
