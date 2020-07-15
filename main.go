package main

//go:generate ./sqlboiler --wipe psql --add-global-variants

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/ystv/web-api/routes"
)

// Version returns web-api's current version
var Version = "dev (0.3.1)"

func main() {
	log.Printf("web-api Version %s", Version)
	err := godotenv.Load(".env.local", ".env") // Load .env file
	if err != nil {
		log.Printf("Failed to load env file %s", err.Error()))
	}
	debug, err := strconv.ParseBool(os.Getenv("debug"))
	if err != nil {
		debug = false
		os.Setenv("debug", "false")
	}
	if debug {
		log.Println("Debug Mode - Disabled auth - pls don't run in production")
	}
	e := routes.Init()

	e.Logger.Fatal(e.Start(":8081"))
}
