package customsettings

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/ystv/web-api/utils"
)

type (
	Repo interface {
		ListCustomSettings(ctx context.Context) ([]CustomSetting, error)
		GetCustomSetting(ctx context.Context, settingID string) (CustomSetting, error)
		AddCustomSetting(ctx context.Context, customSetting CustomSetting) (CustomSetting, error)
		EditCustomSetting(ctx context.Context, settingID string, customSettingEdit CustomSettingEditDTO) (CustomSetting, error)
		DeleteCustomSetting(ctx context.Context, settingID string) error
	}

	CustomSetting struct {
		SettingID string `db:"setting_id" json:"settingID"`
		// Value will be returned in Base64 from PostgreSQL
		Value  interface{} `db:"value" json:"value"`
		Public bool        `db:"public" json:"public"`
	}

	CustomSettingEditDTO struct {
		// Value will be returned in Base64 from PostgreSQL
		Value  string `db:"value" json:"value"`
		Public bool   `db:"public" json:"public"`
	}

	// Store contains our dependency
	Store struct {
		db *sqlx.DB
	}
)

// NewStore creates a new store
func NewStore(db *sqlx.DB) Repo {
	return &Store{db: db}
}

func (s *Store) ListCustomSettings(ctx context.Context) ([]CustomSetting, error) {
	var customSettings []CustomSetting
	//nolint:musttag
	err := s.db.SelectContext(ctx, &customSettings, `
		SELECT setting_id, value, public
		FROM web_api.custom_settings`)
	if err != nil {
		return customSettings, err
	}
	return utils.NonNil(customSettings), nil
}

func (s *Store) GetCustomSetting(ctx context.Context, settingID string) (CustomSetting, error) {
	var customSetting CustomSetting
	//nolint:musttag
	err := s.db.GetContext(ctx, &customSetting, `
		SELECT setting_id, value, public
		FROM web_api.custom_settings
		WHERE setting_id = $1`, settingID)
	if err != nil {
		return customSetting, err
	}

	return customSetting, nil
}

func (s *Store) AddCustomSetting(ctx context.Context, customSetting CustomSetting) (CustomSetting, error) {
	builder := utils.PSQL().Insert("web_api.custom_settings").
		Columns("setting_id", "value", "public").
		Values(customSetting.SettingID, customSetting.Value, customSetting.Public)

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for AddCustomSetting: %w", err))
	}

	res, err := s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return CustomSetting{}, fmt.Errorf("failed to add custom setting: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return CustomSetting{}, fmt.Errorf("failed to add custom setting: %w", err)
	}

	if rows < 1 {
		return CustomSetting{}, fmt.Errorf("failed to add custom setting: invalid rows affected: %d", rows)
	}

	return customSetting, nil
}

func (s *Store) EditCustomSetting(ctx context.Context, settingID string, customSettingEdit CustomSettingEditDTO) (CustomSetting, error) {
	builder := utils.PSQL().Update("web_api.custom_settings").
		SetMap(map[string]interface{}{
			"value":  customSettingEdit.Value,
			"public": customSettingEdit.Public,
		}).
		Where(sq.Eq{"setting_id": settingID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for EditCustomSetting: %w", err))
	}

	res, err := s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return CustomSetting{}, fmt.Errorf("failed to edit custom setting: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return CustomSetting{}, fmt.Errorf("failed to edit custom setting: %w", err)
	}

	if rows < 1 {
		return CustomSetting{}, fmt.Errorf("failed to edit custom setting: invalid rows affected: %d", rows)
	}

	return s.GetCustomSetting(ctx, settingID)
}

func (s *Store) DeleteCustomSetting(ctx context.Context, settingID string) error {
	builder := utils.PSQL().Delete("web_api.custom_settings").
		Where(sq.Eq{"setting_id": settingID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for deleteCustomSetting: %w", err))
	}

	_, err = s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete custom setting: %w", err)
	}

	return nil
}
