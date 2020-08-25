package utils

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// CDN object
var CDN *s3.S3

// InitCDN Initialise connection to CDN
func InitCDN() *s3.S3 {
	endpoint := os.Getenv("cdn_endpoint")
	accessKeyID := os.Getenv("cdn_accessKeyID")
	secretAccessKey := os.Getenv("cdn_secretAccessKey")

	// Configure to use CDN Server

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String("ystv-wales-1"),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession := session.New(s3Config)
	CDN = s3.New(newSession)

	log.Printf("Connected to CDN: %s@%s", accessKeyID, CDN.Endpoint)
	return CDN
}
