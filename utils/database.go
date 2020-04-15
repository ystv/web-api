package utils

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vattle/sqlboiler/boil"
)

func InitDB() *sql.DB {
	err := godotenv.Load() // Load .env file
	if err != nil {
		panic(err)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbURI := fmt.Sprintf("dbname=%s host=%s user=%s password=%s sslmode=disable", dbName, dbHost, username, password) // Build connection string

	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		panic(err)
	}
	err = db.Ping()

	boil.SetDB(db)
	return db
}
