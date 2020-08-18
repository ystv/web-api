package encode

import (
	"github.com/jmoiron/sqlx"
	"github.com/ystv/web-api/utils"
	"gopkg.in/guregu/null.v4"
)

// Preset represents a group of encode formats. A video has a preset applied to it so
// it can generate the video files so a video item.
type Preset struct {
	PresetID    int         `json:"id" db:"id"`
	Name        string      `json:"name" db:"description"`
	Description null.String `json:"description" db:"description"`
	Formats     []Format    `json:"formats"`
}

// PresetList returns all presets
func PresetList() ([]Preset, error) {
	p := []Preset{}
	utils.DB.Select(&p, `SELECT id, name, description
						FROM video.presets;`)
	for i := range p {
		utils.DB.Select(&p[i].Formats,
			`SELECT format.id, name, description, mime_type, mode, width, height, watermarked
			FROM video.encode_formats format
			INNER JOIN video.presets_encode_formats preset ON preset.preset_id = $1;`, p[i].PresetID)
	}
	return p, nil
}

// PresetNew creates a new preset
func PresetNew(p *Preset) error {
	return utils.Transact(utils.DB, func(tx *sqlx.Tx) error {
		presetID := 0
		err := tx.QueryRow("INSERT INTO video.presets(name, description) VALUES ($1, $2) RETURNING id;", p.Name, p.Description).Scan(&presetID)
		if err != nil {
			return err
		}
		stmt, err := tx.Prepare("INSERT INTO video.presets_encode_formats(preset_id, encode_format_id) VALUES ($1, $2);")
		if err != nil {
			return err
		}
		for _, format := range p.Formats {
			_, err := stmt.Exec(presetID, format.FormatID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
