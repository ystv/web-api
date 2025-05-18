-- +goose Up

CREATE SCHEMA IF NOT EXISTS streams;

create sequence streams.incoming_i_id_seq;

create sequence streams.provider_p_id_seq;

create sequence streams.qualities_q_id_seq;

create sequence streams.transcode_map_m_id_seq;

create table streams.incoming
(
    i_id   integer default nextval('streams.incoming_i_id_seq'::regclass) not null
        primary key,
    i_name varchar(256)
);

alter sequence streams.incoming_i_id_seq owned by streams.incoming.i_id;

create table streams.provider
(
    p_id   integer default nextval('streams.provider_p_id_seq'::regclass) not null
        primary key,
    p_user varchar,
    p_host varchar,
    p_key  varchar
);

alter sequence streams.provider_p_id_seq owned by streams.provider.p_id;

create table streams.qualities
(
    q_id   integer default nextval('streams.qualities_q_id_seq'::regclass) not null
        primary key,
    q_name varchar(256),
    q_cmd  varchar(256)
);

alter sequence streams.qualities_q_id_seq owned by streams.qualities.q_id;

create table streams.transcode_map
(
    m_id          integer default nextval('streams.transcode_map_m_id_seq'::regclass) not null
        primary key,
    m_provider_id integer
        references streams.provider
            on update cascade on delete restrict,
    m_quality_id  integer
        references streams.qualities
            on update cascade on delete restrict,
    m_incoming_id integer
        references streams.incoming
            on update cascade on delete restrict
);

alter sequence streams.transcode_map_m_id_seq owned by streams.transcode_map.m_id;

create view "All transcodes"(stream, transcoder, quality) as
SELECT incoming.i_name  AS stream,
       provider.p_host  AS transcoder,
       qualities.q_name AS quality
FROM streams.transcode_map
         JOIN streams.incoming ON transcode_map.m_incoming_id = incoming.i_id
         JOIN streams.qualities ON transcode_map.m_quality_id = qualities.q_id
         JOIN streams.provider ON transcode_map.m_provider_id = provider.p_id
ORDER BY incoming.i_name, qualities.q_id;

-- +goose Down

DROP SCHEMA streams CASCADE;