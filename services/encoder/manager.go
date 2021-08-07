package encoder

import (
	"context"
	"fmt"
)

// Manager subroutine provides a service to manage videos, also
// ensuring consistency of video library.
func Manager() {
	//TODO Make the cool subroutine here
}

// RefreshVideoItem will run CreateEncode() on a VideoItem for any
// encodes missing in the preset.
func (e *Encoder) RefreshVideo(ctx context.Context, videoID int) error {
	// So we will get the video files for a video and the video's preset
	// Check to make sure that there is a source file (we will create renditions based off of it)
	// Check to make sure that there is a preset file set to ensure that encode formats will be created
	v, err := e.getVideoFilesAndPreset(ctx, videoID)
	if err != nil {
		return fmt.Errorf("failed to get video: %w", err)
	}
	if len(v.Files) == 0 {
		return ErrNoVideoFiles
	}
	// We are keeping track of the number of source files since we are ensuring that each
	// video only has one source file, if there is more than one it returns an error
	numOfSrcFiles := 0
	srcFileIdx := 0
	for i, file := range v.Files {
		if file.IsSource {
			numOfSrcFiles += 1
			srcFileIdx = i
		}
	}
	if numOfSrcFiles < 1 {
		return ErrNoSourceFile
	}
	if numOfSrcFiles > 1 {
		return ErrTooManySourceFiles
	}

	if v.PresetID.Valid {
		return ErrNoPreset
	}
	p, err := e.encode.GetPreset(ctx, int(v.PresetID.Int64))
	if err != nil {
		return fmt.Errorf("failed to get preset: %w", err)
	}
	if len(p.Formats) == 0 {
		return ErrNoFormats
	}
	for format := range p.Formats {
		err = e.CreateEncode(ctx, v.Files[srcFileIdx], format)
		if err != nil {
			return fmt.Errorf("failed to create encode fileID=%d format=%d : %w", v.Files[srcFileIdx].FileID, format, err)
		}
	}
	return nil
}

// Refresh will check all existing videoitems to ensure that they
// match their preset, creating new job
func (e *Encoder) Refresh(ctx context.Context) error {

	return nil
}
