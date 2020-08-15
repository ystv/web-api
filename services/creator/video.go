package creator

import (
	"log"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ystv/web-api/services/creator/video"
	"github.com/ystv/web-api/utils"
)

type (
	//IVideoItem defines all creator video interactions
	IVideoItem interface {
		ListVideoItems() ([]video.Item, error)
		FindVideoItems(id int) (video.Item, error)
		CreateVideoItem(item video.Item) (video.Item, error)
		UpdateVideoItem(item video.Item) (video.Item, error)
		DeleteVideoItem(id int) (video.Item, error)
	}
	// PendingUpload represents a uploaded video that didn't have any metadata attached.
	PendingUpload struct {
		ID          int
		Name        string
		Status      string
		Owner       string
		CreatedDate time.Time
	}
	// Preset represents the preset that auto generated the video files from the source material.
	Preset struct {
		ID          int      `json:"id"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Encodes     []Encode `json:"encodes"`
	}
	// Encode represents the each encode of a preset
	Encode struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
		Preset      int    `json:"presetId"`
		Arguements  string `json:"arguements"`
		Watermarked bool   `json:"watermarked"`
	}
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

// ListPendingUploads Returns an array of S3 objects in pending and thier metadata
func ListPendingUploads() ([]PendingUpload, error) {
	videos, err := ListObjects("pending")
	if err != nil {
		log.Printf("Unabled to list objects: %s", err.Error())
		return nil, err
	}
	users := []string{"Rhys", "Rhys but cooler"}
	statuses := []string{"Processing", "Available", "Pending", "Encode error", "Metadata needed", "Locked"}
	rand.Seed(time.Now().Unix())

	var pus []PendingUpload
	for i, video := range videos {
		pu := PendingUpload{ID: i, Name: *video.Key, Owner: users[rand.Intn(len(users))], Status: statuses[rand.Intn(len(statuses))]}
		pus = append(pus, pu)
	}

	return pus, err
}
