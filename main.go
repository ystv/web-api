package main

//go:generate ./sqlboiler --wipe psql --add-global-variants

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/ystv/web-api/routes"
)

func main() {
	err := godotenv.Load(".env.local", ".env") // Load .env file
	if err != nil {
		log.Println(err)
	}
	debug, err := strconv.ParseBool(os.Getenv("debug"))
	if err != nil {
		debug = false
	}
	if debug {
		log.Println("Debug Mode - Disabled auth - pls don't run in production")
	}
	e := routes.Init()
	e.Logger.Fatal(e.Start(":8081"))
}
