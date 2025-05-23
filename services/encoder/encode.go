package encoder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (e *Encoder) getVideoFilesAndPreset(ctx context.Context, videoID int) (VideoItem, error) {
	v := VideoItem{VideoID: videoID}
	//nolint:musttag
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

type EncodeResult struct {
	URI   string
	JobID string
}

// CreateEncode creates an encode item in the message queue.
func (e *Encoder) CreateEncode(ctx context.Context, file VideoFile, formatID int) (EncodeResult, error) {
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
		return EncodeResult{}, fmt.Errorf("failed to get object: %w", err)
	}

	var format EncodeFormat

	err = e.db.GetContext(ctx, &format, `
			SELECT arguments, file_suffix
			FROM video.encode_formats
			WHERE format_id = $1`, formatID)
	if err != nil {
		return EncodeResult{}, err
	}
	if format.Arguments == "" {
		return EncodeResult{}, ErrNoArgs
	}
	if format.FileSuffix == "" {
		format.FileSuffix = strconv.Itoa(formatID)
	}

	// Splitting the URI again, this time on "." so we
	// can apply the file suffix and then put the file
	// extension on after it

	extension := filepath.Ext(key)
	keyWithoutExtension := strings.TrimSuffix(key, extension)

	// Setting the name of the transcoded file
	dstURL := fmt.Sprintf("%s/%s_%s%s", e.conf.ServeBucket, keyWithoutExtension, format.FileSuffix, extension)

	taskVOD := struct {
		SrcURL  string `json:"srcURL"`
		DstArgs string `json:"dstArgs"`
		DstURL  string `json:"dstURL"`
	}{SrcURL: file.URI,
		DstArgs: format.Arguments,
		DstURL:  dstURL}

	reqJSON, err := json.Marshal(taskVOD)
	if err != nil {
		return EncodeResult{}, fmt.Errorf("failed to marshal json: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.conf.VTEndpoint+"/task/video/vod", bytes.NewReader(reqJSON))
	if err != nil {
		return EncodeResult{}, fmt.Errorf("failed to post to vt: %w", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return EncodeResult{}, fmt.Errorf("failed to post to vt: %w", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusCreated:
	case http.StatusUnauthorized:
		return EncodeResult{}, ErrVTFailedToAuthenticate
	default:
		return EncodeResult{}, ErrVTUnknownResponse
	}
	dec := json.NewDecoder(res.Body)

	var task TaskIdentification

	err = dec.Decode(&task)
	if err != nil {
		return EncodeResult{}, fmt.Errorf("failed to decode vt task response: %w", err)
	}

	return EncodeResult{URI: dstURL, JobID: task.TaskID}, nil
}
