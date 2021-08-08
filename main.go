package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/ystv/web-api/controllers/v1/clapper"
	"github.com/ystv/web-api/controllers/v1/creator"
	"github.com/ystv/web-api/controllers/v1/misc"
	"github.com/ystv/web-api/controllers/v1/people"
	"github.com/ystv/web-api/controllers/v1/public"
	"github.com/ystv/web-api/routes"
	"github.com/ystv/web-api/services/encoder"
	"github.com/ystv/web-api/utils"
)

//go:generate swag init -g routes/router.go

// Version returns web-api's current version
var Version = "dev (0.6.4)"

// Commit returns latest commit hash
var Commit = "unknown"

func main() {
	log.Printf("web-api Version %s", Version)
	err := godotenv.Load()            // Load .env file for production
	err = godotenv.Load(".env.local") // Load .env.local for developing
	if err != nil {
		log.Print("Failed to load env file, using global env")
	}

	// Check if debugging
	debug, err := strconv.ParseBool(os.Getenv("WAPI_DEBUG"))
	if err != nil {
		debug = false
		os.Setenv("DEBUG", "false")
	}
	if debug {
		log.Println("Debug Mode - Disabled auth - pls don't run in production")
	}

	// Initialise backend connections
	// Database
	dbConfig := utils.DatabaseConfig{
		Host:     os.Getenv("WAPI_DB_HOST"),
		Port:     os.Getenv("WAPI_DB_PORT"),
		SSLMode:  os.Getenv("WAPI_DB_SSLMODE"),
		Name:     os.Getenv("WAPI_DB_NAME"),
		Username: os.Getenv("WAPI_DB_USER"),
		Password: os.Getenv("WAPI_DB_PASS"),
	}
	db, err := utils.NewDB(dbConfig)
	if err != nil {
		log.Fatalf("failed to start DB: %+v", err)
	}
	log.Printf("Connected to DB: %s@%s", dbConfig.Username, dbConfig.Host)

	// CDN
	cdnConfig := utils.CDNConfig{
		Endpoint:        os.Getenv("WAPI_CDN_ENDPOINT"),
		Region:          os.Getenv("WAPI_CDN_REGION"),
		AccessKeyID:     os.Getenv("WAPI_CDN_ACCESSKEYID"),
		SecretAccessKey: os.Getenv("WAPI_CDN_SECRETACCESSKEY"),
	}
	cdn := utils.NewCDN(cdnConfig)
	log.Printf("Connected to CDN: %s@%s", cdnConfig.AccessKeyID, cdnConfig.Endpoint)

	// Message Queue
	mqConfig := utils.MQConfig{
		Host:     os.Getenv("WAPI_MQ_HOST"),
		Port:     os.Getenv("WAPI_MQ_PORT"),
		Username: os.Getenv("WAPI_MQ_USER"),
		Password: os.Getenv("WAPI_MQ_PASS"),
	}
	mq, err := utils.NewMQ(mqConfig)
	if err != nil {
		log.Fatalf("failed to start mq: %+v", err)
	}
	log.Printf("Connected to MQ: %s@%s", mqConfig.Username, mqConfig.Host)

	// Mail
	// mailPort, err := strconv.Atoi(os.Getenv("WAPI_MAIL_PORT"))
	// if err != nil {
	// 	log.Fatalf("bad mail port: %+v", err)
	// }
	// mailConfig := utils.MailConfig{
	// 	Host:     os.Getenv("WAPI_MAIL_HOST"),
	// 	Port:     mailPort,
	// 	Username: os.Getenv("WAPI_MAIL_USER"),
	// 	Password: os.Getenv("WAPI_MAIL_PASS"),
	// }
	// m, err := utils.NewMailer(mailConfig)
	// if err != nil {
	// 	log.Fatalf("failed to start mailer: %+v", err)
	// }
	// log.Printf("Connected to mail: %s@%s", mailConfig.Username, mailConfig.Host)

	// Messaging
	// utils.InitMessaging()

	bucketConf := struct {
		IngestBucket string
		ServeBucket  string
	}{
		IngestBucket: os.Getenv("WAPI_BUKCET_VOD_INGEST"),
		ServeBucket:  os.Getenv("WAPI_BUCKET_VOD_SERVE"),
	}
	log.Printf("Ingest bucket: %s", bucketConf.IngestBucket)
	log.Printf("Serve bucket: %s", bucketConf.ServeBucket)

	creatorConfig := &creator.Config{
		IngestBucket: bucketConf.IngestBucket,
		ServeBucket:  bucketConf.ServeBucket,
	}

	encoderConfig := &encoder.Config{
		VTEndpoint:  os.Getenv("WAPI_VT_ENDPOINT"),
		ServeBucket: bucketConf.ServeBucket,
	}
	enc := encoder.NewEncoder(db, cdn, mq, encoderConfig)

	routes.New(&routes.NewRouter{
		Version:       Version,
		Commit:        Commit,
		DomainName:    os.Getenv("WAPI_DOMAIN_NAME"),
		Debug:         debug,
		JWTSigningKey: os.Getenv("WAPI_SIGNING_KEY"),
		Clapper:       clapper.NewRepos(db),
		Creator:       creator.NewRepos(db, cdn, enc, creatorConfig),
		Encoder:       enc,
		Misc:          misc.NewRepos(db),
		People:        people.NewRepo(db),
		Public:        public.NewRepos(db),
	}).Start()
}
