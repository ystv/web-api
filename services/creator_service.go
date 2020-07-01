package services

import (
	"log"
	"math/rand"
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
func ListPendingUploads() ([]pendingUpload, error) {
	videos, err := ListObjects("pending")
	if err != nil {
		log.Printf("Unabled to list objects: %s", err.Error())
	}
	users := make([]string, 0)
	users = append(users,
		"Rhys",
		"Rhys but cooler")

	statuses := make([]string, 0)
	statuses = append(statuses,
		"Processing",
		"Available",
		"Pending",
		"Encode error",
		"Metadata needed",
		"Locked")
	rand.Seed(time.Now().Unix())

	var pus []pendingUpload
	for i, video := range videos {
		var pu pendingUpload
		pu.ID = i
		pu.Name = *video.Key
		pu.Owner = users[rand.Intn(len(users))]
		pu.Status = statuses[rand.Intn(len(statuses))]
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

// CreationFind returns the metadata for a given creation
func CreationFind() (*videoitem, error) {
	var creation videoitem
	creation.ID = 1
	creation.Name = "Setup Tour 2020"
	creation.Status = "Available"
	creation.Owner = "Rhys"
	creation.CreatedDate = time.Now()
	creation.Description = "The freshest setup, we got Caspar NRK decided to say screw it, displayed the most fresh NDI out you would have seen"
	var videofile videofile
	videofile.ID = 1
	tmpPreset, _ := PresetFindByID((0))
	videofile.Preset = tmpPreset.Name
	videofile.URI = "cdn.ystv.co.uk/video/1"
	videofile.Status = "Available"
	creation.Files = append(creation.Files, videofile)
	videofile.ID = 2
	videofile.URI = "cdn.ystv.co.uk/video/2"
	tmpPreset, _ = PresetFindByID((2))
	videofile.Preset = tmpPreset.Name
	videofile.Status = "Encode - 4 mins"
	creation.Files = append(creation.Files, videofile)
	return &creation, nil
}

func PresetFindByID(ID int) (*preset, error) {
	var originalVideo preset
	originalVideo.ID = 0
	originalVideo.Name = "Original file"
	originalVideo.Arguements = ""
	originalVideo.Description = "File uploaded to YSTV"
	originalVideo.Watermarked = false
	var hdVideo preset
	hdVideo.ID = 1
	hdVideo.Name = "HD Video"
	hdVideo.Description = "The latest and greatest 720p"
	hdVideo.Arguements = "make it 720i25 pls thanks"
	hdVideo.Watermarked = true
	var unknownVideo preset
	unknownVideo.ID = -1
	unknownVideo.Name = "Unknown video"
	unknownVideo.Description = "We don't know what this file is"
	unknownVideo.Arguements = ""
	unknownVideo.Watermarked = true
	switch ID {
	case 0:
		return &originalVideo, nil
	case 1:
		return &hdVideo, nil
	}
	return &unknownVideo, nil
}

type videoitem struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Status      string      `json:"status"`
	Owner       string      `json:"owner"`
	CreatedDate time.Time   `json:"createdDate"`
	Description string      `json:"description"`
	Files       []videofile `json:"files"`
}

type videofile struct {
	ID     int    `json:"id"`
	URI    string `json:"uri"`
	Preset string `json:"preset"`
	Status string `json:"status"`
}

type preset struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Watermarked bool   `json:"watermarked"`
	Arguements  string `json:"arguements"`
}
