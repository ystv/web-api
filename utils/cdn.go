package utils

import (
	"log"
	"os"

	"github.com/minio/minio-go/v6"
)

// CDN object
var CDN *minio.Client

// InitCDN Initialise connection to CDN
func InitCDN() {
	endpoint := os.Getenv("cdn_endpoint")
	accessKeyID := os.Getenv("cdn_accessKeyID")
	secretAccessKey := os.Getenv("cdn_secretAccessKey")
	useSSL := true

	// Initialise minio client object
	CDN, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to CDN: %s@%s", accessKeyID, CDN.EndpointURL().Host)
}
