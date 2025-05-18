-- +goose Up

CREATE SCHEMA IF NOT EXISTS edit_projects;

CREATE SEQUENCE IF NOT EXISTS edit_projects.cameras_id_seq;

CREATE SEQUENCE IF NOT EXISTS edit_projects.project_id_seq;

CREATE SEQUENCE IF NOT EXISTS edit_projects.sub_project_id_seq;

CREATE TABLE IF NOT EXISTS edit_projects.cameras
(
    id   integer default nextval('edit_projects.cameras_id_seq'::regclass) not null
        primary key,
    name text                                                              not null,
    type text                                                              not null
);

CREATE TABLE IF NOT EXISTS edit_projects.projects
(
    id              integer default nextval('edit_projects.project_id_seq'::regclass) not null
        primary key,
    name            text                                                              not null
        constraint name
            unique,
    has_subprojects boolean                                                           not null,
    director        text                                                              not null
);

CREATE TABLE IF NOT EXISTS edit_projects.sub_projects
(
    project_id     integer                                                               not null
        constraint projects
            references edit_projects.projects
            on update cascade on delete restrict,
    sub_project_id integer default nextval('edit_projects.sub_project_id_seq'::regclass) not null,
    name           text                                                                  not null
        constraint uname
            unique,
    primary key (project_id, sub_project_id)
);

-- +goose Down

DROP SCHEMA edit_projects CASCADE;