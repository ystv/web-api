package encode

import (
	"context"

	"github.com/ystv/web-api/services/creator/types/encode"
)

// ListFormat lists all encode formats
func (s *Store) ListFormat(ctx context.Context) ([]encode.Format, error) {
	e := []encode.Format{}
	err := s.db.SelectContext(ctx, &e, `
		SELECT id, name, description, mime_type, mode, width, height,
		glob_args, src_args, dst_args, watermarked
		FROM video.encode_formats;`)
	if err != nil {
		return nil, err
	}
	return e, nil
}
