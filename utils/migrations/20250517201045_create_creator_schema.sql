-- +goose Up

CREATE SCHEMA IF NOT EXISTS creator;

create table creator.preferences
(
    user_id     integer not null
        references people.users,
    preferences jsonb
);

-- +goose Down

DROP SCHEMA creator CASCADE ;