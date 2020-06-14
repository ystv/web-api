package services

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ystv/web-api/utils"
)

// CreateBucket Creates a new bucket
func CreateBucket(name string, location string) {
	cparams := &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}
	_, err := utils.CDN.CreateBucket(cparams)
	if err != nil {
		log.Printf("Create bucket failed: &v", err)
	}
}

// GenerateUploadURL Creates a signed HTTP POST for a webclient to upload too
func GenerateUploadURL(bucket string, object string) (string, error) {
	// Generates a url which expires in a day.
	expiry := time.Second * 24 * 60 * 60 // 1 day.
	req, _ := utils.CDN.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
	})
	presignedURL, err := req.Presign(expiry)
	//presignedURL, err := utils.CDN.PresignedPutObject(bucket, object, expiry)
	return presignedURL, err
}

// ListObjects Returns an array of ObjectInfo of the input bucket
func ListObjects(bucket string) ([]*s3.Object, error) {
	resp, err := utils.CDN.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})

	return resp.Contents, err
}

// ListPendingUploads Returns an array of S3 objects in pending and thier metadata
func ListPendingUploads() ([]pendingUpload, error) {
	videos, err := ListObjects("pending")
	if err != nil {
		log.Printf("Unabled to list objects: %s", err.Error())
	}

	var pus []pendingUpload
	for i, video := range videos {
		var pu pendingUpload
		pu.ID = i
		pu.Name = *video.Key
		pu.Owner = *video.Owner.DisplayName
		pu.Status = "Processing"
		pus = append(pus, pu)
	}

	return pus, err
}

type pendingUpload struct {
	ID          int
	Name        string
	Status      string
	Owner       string
	CreatedDate time.Time
}
