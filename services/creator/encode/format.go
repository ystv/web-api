package encode

import (
	"context"

	"github.com/ystv/web-api/services/creator"
	"github.com/ystv/web-api/services/creator/types/encode"
)

var _ creator.EncodeRepo = &Store{}

// ListFormat lists all encode formats
func (s *Store) ListFormat(ctx context.Context) ([]encode.Format, error) {
	e := []encode.Format{}
	err := s.db.SelectContext(ctx, &e, `
		SELECT format_id, name, description, mime_type, mode, width, height,
		arguments, file_suffix, watermarked
		FROM video.encode_formats;`)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// NewFormat creates a new format
func (s *Store) NewFormat(ctx context.Context, format encode.Format) (int, error) {
	formatID := 0
	err := s.db.GetContext(ctx, &formatID, `
		INSERT INTO video.encode_formats(name, description, mime_type, mode,
					width, height, arguments, file_suffix, watermarked)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`, format.Name, format.Description,
		format.MimeType, format.Mode, format.Width, format.Height, format.Arguments,
		format.FileSuffix, format.Watermarked)
	if err != nil {
		return 0, err
	}
	return formatID, nil
}

// UpdateFormat will update a format
func (s *Store) UpdateFormat(ctx context.Context, format encode.Format) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE video.encode_formats SET
			name = $1,
			description = $2,
			mime_type = $3,
			mode = $4,
			width = $5,
			height = $6,
			arguments = $7,
			file_suffix = $8,
			watermarked = $9
		WHERE format_id = $10;`, format.Name, format.Description, format.MimeType,
		format.Mode, format.Width, format.Height, format.Arguments,
		format.FileSuffix, format.Watermarked, format.FormatID)
	return err
}

// DeleteFormat will remove a format, videos using this format will need to be removed first
func (s *Store) DeleteFormat(ctx context.Context, formatID int) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM video.encode_formats WHERE format_id = $1`, formatID)
	return err
}
