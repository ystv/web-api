package encoder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (e *Encoder) getVideoFilesAndPreset(ctx context.Context, videoID int) (VideoItem, error) {
	v := VideoItem{}
	v.VideoID = videoID

	err := e.db.GetContext(ctx, &v, `
		SELECT preset_id
		FROM video.items
		WHERE video_id = $1`, videoID)
	if err != nil {
		return v, fmt.Errorf("failed to get video item \"%d\": %w", videoID, err)
	}

	err = e.db.SelectContext(ctx, &v.Files, `
		SELECT file_id, format_id, uri, is_source
		FROM video.files
		WHERE video_id = $1`, videoID)
	if err != nil {
		return v, fmt.Errorf("failed to get video files: %w", err)
	}
	return v, nil
}

// CreateEncode creates an encode item in the message queue.
func (e *Encoder) CreateEncode(ctx context.Context, file VideoFile, formatID int) error {
	// Check video exists
	// Validate encode format
	// Send the job to VT

	URI := strings.Split(file.URI, "/")
	// URI[0] - Bucket
	bucket := URI[0]
	// URI[1] - Key (we have this joiner in the scenario there are multiple slashes in the name)
	key := strings.Join(URI[1:], "")
	_, err := e.cdn.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("failed to get object: %w", err)
	}

	format := EncodeFormat{}
	e.db.GetContext(ctx, &format, `
		SELECT arguments, file_suffix
		FROM video.encode_formats
		WHERE format_id = $1`, formatID)
	if format.Arguments == "" {
		return ErrNoArgs
	}
	if format.FileSuffix == "" {
		format.FileSuffix = fmt.Sprint(formatID)
	}

	// Splitting the URI again, this time on "." so we
	// can apply the file suffix and then put the file
	// extension on after it

	extension := filepath.Ext(key)
	keyWithoutExtension := strings.TrimSuffix(key, extension)

	dstURL := e.conf.ServeBucket + keyWithoutExtension + "_" + format.FileSuffix + extension

	taskVOD := struct {
		SrcURL  string `json:"srcURL"`
		DstArgs string `json:"dstArgs"`
		DstURL  string `json:"dstURL"`
	}{SrcURL: file.URI,
		DstArgs: format.Arguments,
		DstURL:  dstURL}

	reqJSON, err := json.Marshal(taskVOD)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	res, err := e.c.Post(e.conf.VTEndpoint+"/task/video/vod", "application/json", bytes.NewReader(reqJSON))
	if err != nil {
		return fmt.Errorf("failed to post to vt: %w", err)
	}

	defer res.Body.Close()
	switch status := res.StatusCode; {
	case status == http.StatusCreated:
	case status == http.StatusUnauthorized:
		return ErrVTFailedToAuthenticate
	default:
		return ErrVTUnknownResponse
	}
	dec := json.NewDecoder(res.Body)
	task := TaskIdentification{}
	err = dec.Decode(&task)
	if err != nil {
		return fmt.Errorf("failed to decode vt task response: %w", err)
	}
	log.Printf("%+v", task)
	return nil
}
