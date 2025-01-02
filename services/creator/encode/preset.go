package encode

import (
	"context"
	"database/sql"
	"errors"
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
func (s *Store) GetPreset(ctx context.Context, presetID int) (encode.Preset, error) {
	p := encode.Preset{}
	err := s.db.GetContext(ctx, &p, `SELECT preset_id, name, description FROM video.encode_presets;`)
	if err != nil {
		return p, fmt.Errorf("failed to get preset meta: %w", err)
	}
	err = s.db.SelectContext(ctx, &p.Formats,
		`SELECT format.format_id, name, description, mime_type, mode, width, height, watermarked
		FROM video.encode_formats format
		INNER JOIN video.encode_preset_formats preset ON preset.format_id = format.format_id
		WHERE preset.preset_id = $1;`, p.PresetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return p, nil
		}
		return p, fmt.Errorf("failed to get formats: %w", err)
	}
	return p, nil
}

// ListPreset returns all presets
func (s *Store) ListPreset(ctx context.Context) ([]encode.Preset, error) {
	var p []encode.Preset
	//nolint:musttag
	err := s.db.SelectContext(ctx, &p, `SELECT preset_id, name, description
						FROM video.encode_presets;`)
	if err != nil {
		err = fmt.Errorf("failed retrieving meta: %w", err)
		return nil, err
	}
	for i := range p {
		err = s.db.SelectContext(ctx, &p[i].Formats,
			`SELECT format.format_id, name, description, mime_type, mode, width, height, watermarked
			FROM video.encode_formats format
			INNER JOIN video.encode_preset_formats preset ON preset.format_id = format.format_id
			WHERE preset.preset_id = $1;`, p[i].PresetID)
		if err != nil {
			err = fmt.Errorf("failed retrieving formats: %w", err)
			return nil, err
		}
	}
	return p, nil
}

// NewPreset creates a new preset
func (s *Store) NewPreset(ctx context.Context, p encode.Preset) (int, error) {
	presetID := 0
	err := utils.Transact(s.db, func(tx *sqlx.Tx) error {
		err := tx.GetContext(ctx, &presetID, "INSERT INTO video.encode_presets(name, description) VALUES ($1, $2) RETURNING preset_id;", p.Name, p.Description)
		if err != nil {
			return fmt.Errorf("failed to insert preset meta: %w", err)
		}

		// When they don't attach any formats
		if len(p.Formats) == 0 {
			return nil
		}

		stmt, err := tx.PrepareContext(ctx, "INSERT INTO video.encode_preset_formats(preset_id, format_id) VALUES ($1, $2);")
		if err != nil {
			return fmt.Errorf("failed to prepare statement to insert formats: %w", err)
		}
		defer stmt.Close()

		for _, format := range p.Formats {
			_, err := stmt.ExecContext(ctx, presetID, format.FormatID)
			if err != nil {
				return fmt.Errorf("failed to insert link between preset and formats: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create new preset: %w", err)
	}
	return presetID, nil
}

// UpdatePreset updates an existing preset
func (s *Store) UpdatePreset(ctx context.Context, p encode.Preset) error {
	return utils.Transact(s.db, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE video.encode_presets SET name = $1, description = $2
							WHERE preset_id = $3;`, p.Name, p.Description, p.PresetID)
		if err != nil {
			return fmt.Errorf("failed to update preset meta: %w", err)
		}

		// Deleting old associated encode formats
		_, err = tx.ExecContext(ctx, `DELETE FROM video.encode_preset_formats
						WHERE preset_id = $1`, p.PresetID)
		if err != nil {
			return fmt.Errorf("failed to delete old format links: %w", err)
		}

		// When they don't attach any formats
		if len(p.Formats) == 0 {
			return nil
		}

		// Insert new formats
		stmt, err := tx.PrepareContext(ctx, "INSERT INTO video.encode_preset_formats(preset_id, format_id) VALUES ($1, $2);")
		if err != nil {
			return fmt.Errorf("failed to prepare statement to insert formats: %w", err)
		}
		defer stmt.Close()

		for _, format := range p.Formats {
			_, err := stmt.ExecContext(ctx, p.PresetID, format.FormatID)
			if err != nil {
				return fmt.Errorf("failed to insert link between preset and formats: %w", err)
			}
		}

		return nil
	})
}

// DeletePreset deletes a preset, this won't affected any formats that are part of the preset
func (s *Store) DeletePreset(ctx context.Context, presetID int) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM video.encode_presets WHERE preset_id = $1`, presetID)
	return err
}
