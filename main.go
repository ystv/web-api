package main

//go:generate ./sqlboiler --wipe psql --add-global-variants

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/ystv/web-api/routes"
)

func main() {
	err := godotenv.Load() // Load .env file
	if err != nil {
		log.Println("No .env file present, using global env")
	}
	e := routes.Init()
	e.Logger.Fatal(e.Start(":8080"))
}
