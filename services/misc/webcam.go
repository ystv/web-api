package misc

import (
	"context"
	"fmt"
)

type (
	// Webcam represents a watchable webcam
	Webcam struct {
		CameraID int    `db:"camera_id" json:"id"`
		Name     string `db:"name" json:"name"`
		URL      string `db:"url" json:"-"`
		File     string `db:"file" json:"file"`
		MIMEType string `db:"mime_type" json:"mimeType"`
	}
	// AdminWebcam represents extra options to configure the webcam
	AdminWebcam struct {
		Webcam
		Enabled    bool   `db:"enabled" json:"enabled"`
		Permission string `db:"permission" json:"permission"`
	}
)

// Here for validation to ensure we are meeting the interface
var _ WebcamRepo = &Store{}

// ListWebcams returns all webcams a user can access
func (m *Store) ListWebcams(ctx context.Context, permissions []string) ([]Webcam, error) {
	webcams := []AdminWebcam{}
	publicWebcams := []Webcam{}
	// Fetch all enabled webcams from DB
	err := m.db.SelectContext(ctx, &webcams,
		`SELECT	camera_id, name, file, mime_type, permission
		FROM misc.webcams
		WHERE ENABLED;`)
	if err != nil {
		return publicWebcams, fmt.Errorf("failed to select webcams: %w", err)
	}

	// Check if user has permission to view it
	publicWebcam := Webcam{}
	isAuthorized := false
	for _, webcam := range webcams {
		isAuthorized = false
		for _, perm := range permissions {
			if webcam.Permission == perm || webcam.Permission == "" {
				publicWebcam = Webcam{
					webcam.CameraID,
					webcam.Name,
					webcam.URL,
					webcam.File,
					webcam.MIMEType,
				}
				isAuthorized = true
			}
		}
		if isAuthorized {
			publicWebcams = append(publicWebcams, publicWebcam)
		}
	}
	return publicWebcams, nil
}

// GetWebcam returns a single webcam
func (m *Store) GetWebcam(ctx context.Context, cameraID int, permissions []string) (Webcam, error) {
	webcam := AdminWebcam{}
	publicWebcam := Webcam{}
	err := m.db.GetContext(ctx, &webcam,
		`SELECT	camera_id, name, url, file, mime_type, permission_id
		FROM misc.webcams
		WHERE ENABLED AND
		camera_id = $1;`, cameraID)
	if err != nil {
		err = fmt.Errorf("failed to select webcams: %w", err)
		return publicWebcam, err
	}

	// Check if user has permission to view it
	for _, perm := range permissions {
		if webcam.Permission == perm || webcam.Permission == "" {
			publicWebcam = Webcam{
				webcam.CameraID,
				webcam.Name,
				webcam.URL,
				webcam.File,
				webcam.MIMEType,
			}
		}
	}
	return publicWebcam, nil
}
