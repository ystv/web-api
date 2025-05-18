-- +goose Up

CREATE TABLE web_api.custom_settings
(
    setting_id    varchar                  not null
        constraint custom_settings_pk
            primary key,
    value        varchar                  not null,
    public boolean default false
);

-- +goose Down

DROP TABLE web_api.custom_settings;