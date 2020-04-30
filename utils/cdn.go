package utils

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v6"
)

var CDN *minio.Client

// InitCDN Initialise connection to CDN
func InitCDN() {
	err := godotenv.Load() // Load .env file
	if err != nil {
		panic(err)
	}
	endpoint := os.Getenv("cdn_endpoint")
	accessKeyID := os.Getenv("cdn_accessKeyID")
	secretAccessKey := os.Getenv("cdn_secretAccessKey")
	useSSL := true

	// Initialise minio client object
	CDN, err = minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		panic(err)
	}
}
