package public

import "context"

type CustomSetting struct {
	SettingID string `db:"setting_id" json:"settingID"`
	// Value will be returned in Base64 from PostgreSQL
	Value interface{} `db:"value" json:"value"`
}

func (s *Store) GetCustomSettingPublic(ctx context.Context, settingID string) (CustomSetting, error) {
	var customSetting CustomSetting
	err := s.db.GetContext(ctx, &customSetting, `
		SELECT setting_id, value
		FROM web_api.custom_settings
		WHERE setting_id = $1 AND public = true;`, settingID)
	if err != nil {
		return customSetting, err
	}

	return customSetting, nil
}
