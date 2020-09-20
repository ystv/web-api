package misc

import (
	"context"
	"fmt"
)

type (
	// Webcam represents a watchable webcam
	Webcam struct {
		CameraID int    `db:"camera_id" json:"id"`
		URL      string `db:"url" json:"url"`
	}
	// AdminWebcam represents extra options to configure the webcam
	AdminWebcam struct {
		Webcam
		Enabled      bool `db:"enabled" json:"enabled"`
		PermissionID int  `db:"permission_id" json:"permissionID"`
	}
)

// Here for validation to ensure we are meeting the interface
var _ WebcamRepo = &Store{}

// ListWebcams returns all webcams a user can access
func (m *Store) ListWebcams(ctx context.Context, permissionIDs []int) ([]Webcam, error) {
	w := []AdminWebcam{}
	publicWebcams := []Webcam{}
	// Fetch all enabled webcams from DB
	err := m.db.SelectContext(ctx, &w,
		`SELECT	camera_id, url, permission_id
		FROM misc.webcams
		WHERE ENABLED;`)
	if err != nil {
		err = fmt.Errorf("failed to select webcams: %w", err)
		return publicWebcams, err
	}

	// Check if user has permission to view it
	publicWebcam := Webcam{}
	for _, webcam := range w {
		for _, id := range permissionIDs {
			if id == webcam.PermissionID {
				publicWebcam = Webcam{
					webcam.CameraID,
					webcam.URL,
				}
				publicWebcams = append(publicWebcams, publicWebcam)
			}
		}

	}

	return publicWebcams, nil
}

// GetWebcam returns a single webcam
func (m *Store) GetWebcam(ctx context.Context, cameraID int, permissionIDs []int) (Webcam, error) {
	w := AdminWebcam{}
	publicWebcam := Webcam{}
	err := m.db.GetContext(ctx, &w,
		`SELECT	camera_id, url, permission_id
		FROM misc.webcams
		WHERE ENABLED;`)
	if err != nil {
		err = fmt.Errorf("failed to select webcams: %w", err)
		return publicWebcam, err
	}

	// Check if user has permission to view it
	for _, id := range permissionIDs {
		if id == w.PermissionID {
			publicWebcam = Webcam{
				w.CameraID,
				w.URL,
			}
		}
	}
	return publicWebcam, nil
}
