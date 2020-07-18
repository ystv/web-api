package creator

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ystv/web-api/storage"
	"github.com/ystv/web-api/utils"
)

type (
	//IVideoItem defines all creator video interactions
	IVideoItem interface {
		ListVideoItems() ([]VideoItem, error)
		FindVideoItems(id int) (VideoItem, error)
		CreateVideoItem(item VideoItem) (VideoItem, error)
		UpdateVideoItem(item VideoItem) (VideoItem, error)
		DeleteVideoItem(id int) (VideoItem, error)
	}
	// PendingUpload represents a uploaded video that didn't have any metadata attached.
	PendingUpload struct {
		ID          int
		Name        string
		Status      string
		Owner       string
		CreatedDate time.Time
	}
	// VideoMeta represents basic information about the videoitem used for listing.
	VideoMeta struct {
		ID            int       `json:"id"`
		Name          string    `json:"name"`
		Description   string    `json:"description"`
		Status        string    `json:"status"`
		Owner         string    `json:"owner"`
		BroadcastDate time.Time `json:"broadcastDate"`
		Views         int       `json:"views"`
		Duration      int       `json:"duration"`
		Preset        string    `json:"preset"`
	}
	// VideoItem represents the basic in-depth information including videofiles.
	VideoItem struct {
		ID          int         `json:"id"`
		Name        string      `json:"name"`
		Status      string      `json:"status"`
		Owner       string      `json:"owner"`
		CreatedDate time.Time   `json:"createdDate"`
		Description string      `json:"description"`
		Duration    int         `json:"duration"`
		Preset      string      `json:"preset"`
		Views       int         `json:"views"`
		Files       []VideoFile `json:"files"`
	}
	// VideoFile represents each file that a video item has stored.
	VideoFile struct {
		ID     int    `json:"id"`
		URI    string `json:"uri"`
		Preset string `json:"preset"`
		Status string `json:"status"`
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

// VideoItemFind returns the metadata for a given creation
func VideoItemFind(ctx context.Context, id int) (*storage.SQLVideoItem, error) {
	return storage.VideoItem(ctx, id)

	// files, err := json.Marshal(v.Files)
	// log.Printf("Marshal: %+v", err)
	// log.Printf(string(files))
	// files2 := DBVideoFile{}
	// err = json.Unmarshal(files, &files2)
	// log.Printf("Unmarshal: %+v", err)
	// log.Printf("new array: %+v", files2)
	// v2 := DBVideoItem2{ID: v.ID, Name: v.Name, Files: files2}
	// return &v2, nil
	//creation := VideoItem{
	//	ID:          1,
	//	Name:        "Setup Tour 2020",
	//	Status:      "Available",
	//	Owner:       "Rhys",
	//	CreatedDate: time.Now(),
	//	Description: "Big video description",
	//	Duration:    300,
	//	Views:       56,
	//	Files: []VideoFile{{
	//		ID:     1,
	//		Preset: "Original master",
	//		Status: "Internal",
	//		URI:    "cdn.ystv.co.uk/videos/1/1",
	//	}, {
	//		ID:     2,
	//		Preset: "FHD Video",
	//		Status: "Public",
	//		URI:    "cdn.ystv.co.uk/videos/1/2",
	//	}, {
	//		ID:     3,
	//		Preset: "HD Video",
	//		Status: "Processing",
	//		URI:    "",
	//	}, {
	//		ID:     4,
	//		Preset: "English Subtitles",
	//		Status: "Public",
	//		URI:    "cdn.ystv.co.uk/videos/1/4",
	//	}, {
	//		ID:     5,
	//		Preset: "Thumbnails",
	//		Status: "Internal",
	//		URI:    "cdn.ystv.co.uk/videos/1/5",
	//	},
	//	},
	//}
	//return &creation, nil
}

// PresetFind a preset from it's ID
func PresetFind(ID int) (*Preset, error) {
	originalVideo := Preset{ID: 0, Name: "Original file", Description: "File uploaded to YSTV"}
	hdVideo := Preset{ID: 1, Name: "HD Video", Description: "The latest and greatest 720p"}
	unknownVideo := Preset{ID: -1, Name: "Unknown video", Description: "We don't know what this file is"}
	switch ID {
	case 0:
		return &originalVideo, nil
	case 1:
		return &hdVideo, nil
	}
	return &unknownVideo, nil
}
