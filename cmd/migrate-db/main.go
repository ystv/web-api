package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"

	"github.com/ystv/web-api/utils"
	"github.com/ystv/web-api/utils/migrations"
)

func main() {
	// Load environment
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("failed to load global env file")
	} // Load .env file for production
	err = godotenv.Overload(".env.local") // Load .env.local for developing
	if err != nil {
		log.Println("failed to load env file, using global env")
	}

	downOne := flag.Bool("down_one", false, "undo the last migration instead of upgrading - only use for development!")
	flag.Parse()

	host := os.Getenv("WAPI_DB_HOST")

	if host == "" {
		log.Fatalf("database host not set")
	}
	dbConfig := utils.DatabaseConfig{
		Host:     host,
		Port:     os.Getenv("WAPI_DB_PORT"),
		SSLMode:  os.Getenv("WAPI_DB_SSLMODE"),
		Name:     os.Getenv("WAPI_DB_NAME"),
		Username: os.Getenv("WAPI_DB_USER"),
		Password: os.Getenv("WAPI_DB_PASS"),
	}
	database, err := utils.NewDB(dbConfig)

	goose.SetBaseFS(migrations.Migrations)
	if err = goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}

	if *downOne {
		if err = goose.Down(database.DB, "."); err != nil {
			log.Fatalf("unable to downgrade: %v", err)
		}
		return
	}

	if err = goose.Up(database.DB, "."); err != nil {
		log.Fatalf("unable to run migrations: %v", err)
	}

	log.Println("migrations ran successfully")
}
