package main

//go:generate ./sqlboiler --wipe psql --add-global-variants

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/ystv/web-api/routes"
	"github.com/ystv/web-api/utils"
)

// Version returns web-api's current version
var Version = "dev (0.5.3)"

// Commit returns latest commit hash
var Commit = "unknown"

func main() {
	log.Printf("web-api Version %s", Version)
	err := godotenv.Load()            // Load .env file for production
	err = godotenv.Load(".env.local") // Load .env.local for developing
	if err != nil {
		log.Print("Failed to load env file, using global env")
	}
	debug, err := strconv.ParseBool(os.Getenv("debug"))
	if err != nil {
		debug = false
		os.Setenv("debug", "false")
	}
	if debug {
		log.Println("Debug Mode - Disabled auth - pls don't run in production")
	}
	utils.InitDB()
	utils.InitCDN()
	// utils.InitMessaging()

	e := routes.Init(Version, Commit)

	e.Logger.Fatal(e.Start(":8081"))
}
