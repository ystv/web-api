package encoder

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ystv/web-api/services/creator"
	"github.com/ystv/web-api/utils"
)

// CreateBucket Creates a new bucket
func CreateBucket(name string, location string) {
	cparams := &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}
	_, err := utils.CDN.CreateBucket(cparams)
	if err != nil {
		log.Printf("Create bucket failed: %v", err)
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

// Manager subroutine provides a service to manage videos, also
// ensuring consistency of video library.

var MessagingClient utils.IMessagingClient

func Manager() {
	//TODO Make the cool subroutine here
}

// CreateEncode creates an encode item in the message queue.
func CreateEncode(v creator.VideoFile, e creator.Encode) error {
	return nil
}

func ListEncodesFromPreset(p creator.Preset) ([]creator.Encode, error) {
	return nil, nil
}

// RefreshVideoItem will run CreateEncode() on a VideoItem for any
// encodes missing in the preset.
func RefreshVideoItem(v creator.VideoItem) error {
	return nil
}

// Refresh will check all existing videoitems to ensure that they
// match their preset, creating new job
func Refresh() error {
	return nil
}