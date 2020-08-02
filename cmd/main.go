package main

//go:generate ./sqlboiler --wipe psql --add-global-variants

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/ystv/web-api/routes"
	"github.com/ystv/web-api/services/creator"
	"github.com/ystv/web-api/services/creator/playlist"
	"github.com/ystv/web-api/utils"
	"gopkg.in/guregu/null.v4"
)

// Version returns web-api's current version
var Version = "dev (0.5.1)"

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

	p := playlist.Playlist{
		Meta: playlist.Meta{
			ID:          3,
			Name:        "Funny videos!",
			Description: null.StringFrom("The most epic videos where each submitted video gets £200?!"),
			Status:      "internal",
			CreatedAt:   time.Now(),
			CreatedBy:   1,
		},
	}
	v, _ := creator.VideoMetaList(context.Background())
	single := (*v)[15]
	log.Print(playlist.AddVideo(p.Meta, &single))
	item, _ := playlist.Get(3)
	log.Print(item)

	e := routes.Init(Version, Commit)

	e.Logger.Fatal(e.Start(":8081"))
}
