package encode

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ystv/web-api/services/creator"
	"github.com/ystv/web-api/services/creator/types/encode"
	"github.com/ystv/web-api/utils"
)

var _ creator.EncodeRepo = &Store{}

// Store contains our dependency
type Store struct {
	db *sqlx.DB
}

// NewStore creates our data store
func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

// GetPreset returns a preset by ID
func (s *Store) GetPreset(ctx context.Context, presetID int) (*encode.Preset, error) {
	p := encode.Preset{}
	err := s.db.GetContext(ctx, &p, `SELECT id, name, description
						FROM video.presets;`)
	if err != nil {
		return nil, err
	}
	err = s.db.SelectContext(ctx, &p.Formats,
		`SELECT format.id, name, description, mime_type, mode, width, height, watermarked
		FROM video.encode_formats format
		INNER JOIN video.presets_encode_formats preset ON preset.encode_format_id = format.id
		WHERE preset.preset_id = $1;`, p.PresetID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListPreset returns all presets
func (s *Store) ListPreset(ctx context.Context) ([]encode.Preset, error) {
	p := []encode.Preset{}
	err := s.db.SelectContext(ctx, &p, `SELECT id, name, description
						FROM video.presets;`)
	if err != nil {
		err = fmt.Errorf("failed retrieving meta: %w", err)
		return nil, err
	}
	for i := range p {
		err = s.db.SelectContext(ctx, &p[i].Formats,
			`SELECT format.id, name, description, mime_type, mode, width, height, watermarked
			FROM video.encode_formats format
			INNER JOIN video.presets_encode_formats preset ON preset.encode_format_id = format.id
			WHERE preset.preset_id = $1;`, p[i].PresetID)
		if err != nil {
			err = fmt.Errorf("failed retrieving formats: %w", err)
			return nil, err
		}
	}
	return p, nil
}

// NewPreset creates a new preset
func (s *Store) NewPreset(ctx context.Context, p *encode.Preset) (int, error) {
	return 0, utils.Transact(s.db, func(tx *sqlx.Tx) error {
		presetID := 0
		err := tx.QueryRowContext(ctx, "INSERT INTO video.presets(name, description) VALUES ($1, $2) RETURNING id;", p.Name, p.Description).Scan(&presetID)
		if err != nil {
			err = fmt.Errorf("failed to insert preset meta: %w", err)
			return err
		}
		// When they don't attach any formats
		if len(p.Formats) == 0 {
			return nil
		}
		stmt, err := tx.PrepareContext(ctx, "INSERT INTO video.presets_encode_formats(preset_id, encode_format_id) VALUES ($1, $2);")
		if err != nil {
			err = fmt.Errorf("failed to prepare statement to insert formats: %w", err)
			return err
		}
		for _, format := range p.Formats {
			_, err := stmt.ExecContext(ctx, presetID, format.FormatID)
			if err != nil {
				err = fmt.Errorf("failed to inset link between preset and formats: %w", err)
				return err
			}
		}
		return nil
	}) // TODO return preset ID
}

// UpdatePreset updates an existing preset
func (s *Store) UpdatePreset(ctx context.Context, p *encode.Preset) error {
	return utils.Transact(s.db, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE video.presets SET name = $1, description = $2
							WHERE id = $3;`, p.Name, p.Description, p.PresetID)
		if err != nil {
			err = fmt.Errorf("failed to update preset meta: %w", err)
			return err
		}
		// Deleting old associated encode formats
		_, err = tx.ExecContext(ctx, `DELETE FROM video.presets_encode_formats
						WHERE preset_id = $1`, p.PresetID)
		if err != nil {
			err = fmt.Errorf("failed to delete old format links: %w", err)
			return err
		}
		// When they don't attach any formats
		if len(p.Formats) == 0 {
			return nil
		}
		// Insert new formats
		stmt, err := tx.PrepareContext(ctx, "INSERT INTO video.presets_encode_formats(preset_id, encode_format_id) VALUES ($1, $2);")
		if err != nil {
			err = fmt.Errorf("failed to prepare statement to insert formats: %w", err)
			return err
		}
		for _, format := range p.Formats {
			_, err := stmt.ExecContext(ctx, p.PresetID, format.FormatID)
			if err != nil {
				err = fmt.Errorf("failed to insert link between preset and formats: %w", err)
				return err
			}
		}
		return nil
	})
}
