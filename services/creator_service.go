package services

import (
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v6"
	"github.com/ystv/web-api/utils"
)

// CreateBucket Creates a new bucket
func CreateBucket(name string, location string) {
	err := utils.CDN.MakeBucket(name, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := utils.CDN.BucketExists(name)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", name)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", name)
	}
}

// GenerateUploadURL Creates a signed HTTP POST for a webclient to upload too
func GenerateUploadURL(bucket string, object string) (*url.URL, error) {
	// Generates a url which expires in a day.
	expiry := time.Second * 24 * 60 * 60 // 1 day.
	presignedURL, err := utils.CDN.PresignedPutObject(bucket, object, expiry)
	return presignedURL, err
}

// ListObjects Returns an array of ObjectInfo of the input bucket
func ListObjects(bucket string) []minio.ObjectInfo {
	var objects []minio.ObjectInfo
	// Create a done channel to control function go routine.
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return
	defer close(doneCh)

	objectCh := utils.CDN.ListObjectsV2(bucket, "/", true, doneCh)
	for object := range objectCh {
		if object.Err != nil {
			log.Printf("Couldn't list object: %v", object.Err)
		}
		objects = append(objects, object)
	}
	return objects
}
