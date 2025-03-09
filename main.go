package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/ystv/web-api/controllers/v1/clapper"
	"github.com/ystv/web-api/controllers/v1/creator"
	encoderPackage "github.com/ystv/web-api/controllers/v1/encoder"
	"github.com/ystv/web-api/controllers/v1/misc"
	"github.com/ystv/web-api/controllers/v1/people"
	"github.com/ystv/web-api/controllers/v1/public"
	"github.com/ystv/web-api/controllers/v1/stream"
	"github.com/ystv/web-api/services/encoder"
	"github.com/ystv/web-api/utils"
)

//go:generate swag init -o swagger/ --overridesFile swagger/.swaggo

// Version returns web-api's current version
var Version = "dev (0.8.0)"

// Commit returns latest commit hash
var Commit = "unknown"

// @title YSTV Web API
// @version 1.0
// @description This is YSTV's API for public and backend sites
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.url https://comp.ystv.co.uk
// @contact.email computing@ystv.co.uk
//
// @license.name GPL-3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.en.html
//
// @host api.ystv.co.uk
func main() {
	var local, global bool

	var err error
	err = godotenv.Load(".env") // Load .env
	global = err == nil

	err = godotenv.Overload(".env.local") // Load .env.local
	local = err == nil

	dbHost := os.Getenv("WAPI_DB_HOST")

	if !local && !global && dbHost == "" {
		log.Fatal("unable to find env files and no env variables have been supplied")
	}
	//nolint:gocritic
	if !local && !global {
		log.Println("using env variables")
	} else if local && global {
		log.Println("using global and local env files")
	} else if !local {
		log.Println("using global env file")
	} else {
		log.Println("using local env file")
	}

	log.Printf("web-api version: %s, commit: %s\n", Version, Commit)

	// Check if debugging
	debug, err := strconv.ParseBool(os.Getenv("WAPI_DEBUG"))
	if err != nil {
		debug = false
		err = os.Setenv("DEBUG", "false")
		if err != nil {
			log.Printf("failed to set env: %v", err)
		}
	}
	if debug {
		log.Println("debug Mode - disabled auth - don't run in production")
	}

	// Initialise backend connections
	// Database
	dbConfig := utils.DatabaseConfig{
		Host:     dbHost,
		Port:     os.Getenv("WAPI_DB_PORT"),
		SSLMode:  os.Getenv("WAPI_DB_SSLMODE"),
		Name:     os.Getenv("WAPI_DB_NAME"),
		Username: os.Getenv("WAPI_DB_USER"),
		Password: os.Getenv("WAPI_DB_PASS"),
	}
	db, err := utils.NewDB(dbConfig)
	if err != nil {
		log.Fatalf("failed to connect db: %+v", err)
	}
	log.Printf("connected to db: %s", dbConfig.Host)

	// CDN
	cdnConfig := utils.CDNConfig{
		Endpoint:        os.Getenv("WAPI_CDN_ENDPOINT"),
		Region:          os.Getenv("WAPI_CDN_REGION"),
		AccessKeyID:     os.Getenv("WAPI_CDN_ACCESSKEYID"),
		SecretAccessKey: os.Getenv("WAPI_CDN_SECRETACCESSKEY"),
	}
	cdn, err := utils.NewCDN(cdnConfig)
	if err != nil {
		log.Fatalf("unable to connect to cdn: %v", err)
	}
	log.Printf("connected to cdn: %s", cdnConfig.Endpoint)

	bucketConf := struct {
		IngestBucket string
		ServeBucket  string
	}{
		IngestBucket: os.Getenv("WAPI_BUCKET_VOD_INGEST"),
		ServeBucket:  os.Getenv("WAPI_BUCKET_VOD_SERVE"),
	}

	jwtCookieName := os.Getenv("WAUTH_JWT_COOKIE_NAME")
	if jwtCookieName == "" {
		jwtCookieName = "wauth_jwt"
	}

	access := utils.NewAccesser(utils.Config{
		AccessCookieName: jwtCookieName,
		SigningKey:       []byte(os.Getenv("WAPI_SIGNING_KEY")),
	})

	creatorConfig := &creator.Config{
		IngestBucket: bucketConf.IngestBucket,
		ServeBucket:  bucketConf.ServeBucket,
	}

	encoderConfig := &encoder.Config{
		VTEndpoint:  os.Getenv("WAPI_VT_ENDPOINT"),
		ServeBucket: bucketConf.ServeBucket,
	}
	enc := encoder.NewEncoder(db, cdn, encoderConfig)

	New(&NewRouter{
		Version:    Version,
		Commit:     Commit,
		DomainName: os.Getenv("WAPI_DOMAIN_NAME"),
		Debug:      debug,
		Access:     access,
		Clapper:    clapper.NewRepos(db, access),
		Creator:    creator.NewRepos(db, cdn, enc, access, creatorConfig, cdnConfig.Endpoint),
		Encoder:    encoderPackage.NewEncoderController(enc, access),
		Misc:       misc.NewRepos(db, access),
		People:     people.NewRepos(db, cdn, access, cdnConfig.Endpoint),
		Public:     public.NewRepos(db),
		Stream:     stream.NewRepos(db),
	}).Start()
}
