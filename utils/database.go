package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/vattle/sqlboiler/boil"
)

// DB object
var DB *sql.DB

// InitDB Initialises the connection to the database
func InitDB() {
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")

	dbURI := fmt.Sprintf("dbname=%s host=%s user=%s password=%s port=%s sslmode=disable", dbName, dbHost, username, password, dbPort) // Build connection string

	DB, err := sql.Open("postgres", dbURI)
	if err != nil {
		panic(err)
	}
	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	boil.SetDB(DB)

	ret, err := strconv.ParseBool(os.Getenv("debug"))
	if err != nil {
		panic(err)
	}
	boil.DebugMode = ret
	log.Printf("Connected to DB: %s@%s", dbName, dbHost)
}
