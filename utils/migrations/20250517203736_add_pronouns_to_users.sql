-- +goose Up

alter table people.users
    add pronouns text;

-- +goose Down

alter table people.users
    drop column pronouns;