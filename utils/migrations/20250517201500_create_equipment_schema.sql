-- +goose Up

CREATE SCHEMA IF NOT EXISTS equipment;

CREATE SEQUENCE IF NOT EXISTS equipment.categories_cat_id_seq;

CREATE SEQUENCE IF NOT EXISTS equipment.equipment_id_seq;

CREATE SEQUENCE IF NOT EXISTS equipment.hire_categories_cat_id_seq;

CREATE SEQUENCE IF NOT EXISTS equipment.hires_description_seq;

CREATE SEQUENCE IF NOT EXISTS equipment.hires_id_seq;

CREATE SEQUENCE IF NOT EXISTS equipment.logs_id_seq;

CREATE TABLE IF NOT EXISTS equipment.categories
(
    category    varchar(30)                                                          not null
        constraint category_unique
            unique,
    description text,
    cat_id      integer default nextval('equipment.categories_cat_id_seq'::regclass) not null
        primary key
);

COMMENT ON COLUMN equipment.categories.category is 'Name for the category';

COMMENT ON COLUMN equipment.categories.description is 'More information';

CREATE TABLE IF NOT EXISTS equipment.equipment
(
    name          text                                                               not null,
    description   text,
    product_num   text,
    serial_num    text,
    source        text,
    cost          money,
    date_acquired date,
    hire          boolean    default false                                           not null,
    disposed      boolean    default false                                           not null,
    asset_num     numeric(4) default NULL::numeric
        unique,
    date_disposed date,
    id            integer    default nextval('equipment.equipment_id_seq'::regclass) not null
        primary key,
    ystv_name     text,
    category_id   integer
        references equipment.categories
            on update cascade
);

COMMENT ON COLUMN equipment.equipment.name is 'Product name of the item';

COMMENT ON COLUMN equipment.equipment.description is 'More information if required';

COMMENT ON COLUMN equipment.equipment.product_num is 'Generic product number for the item';

COMMENT ON COLUMN equipment.equipment.serial_num is 'Specific serial number for the item';

COMMENT ON COLUMN equipment.equipment.source is 'Where it was acquired from';

COMMENT ON COLUMN equipment.equipment.cost is 'How much was it';

COMMENT ON COLUMN equipment.equipment.date_acquired is 'When we got it';

COMMENT ON COLUMN equipment.equipment.hire is 'Is it for hire';

COMMENT ON COLUMN equipment.equipment.disposed is 'Have we got rid of it';

COMMENT ON COLUMN equipment.equipment.asset_num is 'Asset number, as is printed on the tag';

COMMENT ON COLUMN equipment.equipment.date_disposed is 'When we got rid';

COMMENT ON COLUMN equipment.equipment.ystv_name is 'What YSTV has decided to call this specific item';

COMMENT ON COLUMN equipment.equipment.category_id is 'What things are';

CREATE TABLE IF NOT EXISTS equipment.hire_categories
(
    category    text                                                                      not null
        constraint hire_category_unique
            unique,
    description text,
    cat_id      integer default nextval('equipment.hire_categories_cat_id_seq'::regclass) not null
        primary key
);

CREATE TABLE IF NOT EXISTS equipment.hires
(
    item             text                                                        not null
        unique,
    hire_cost        money,
    hire_description text,
    category         text
        references equipment.hire_categories (category)
            on update cascade,
    id               integer default nextval('equipment.hires_id_seq'::regclass) not null
        primary key,
    category_id      integer default 9                                           not null
        references equipment.hire_categories
            on update cascade
);

COMMENT ON COLUMN equipment.hires.category is 'Deprecated. May be dropped without warning';

CREATE TABLE IF NOT EXISTS equipment.logs
(
    id          integer                  default nextval('equipment.logs_id_seq'::regclass) not null
        primary key,
    log         text                                                                        not null,
    posted_by   integer
        references people.users
            on update cascade,
    posted_date timestamp with time zone default now()                                      not null,
    item_id     integer                                                                     not null
        references equipment.equipment
            on update cascade
);

COMMENT ON COLUMN equipment.logs.posted_by is 'Member id of creator of log';

CREATE VIEW "YUSU List" ("Quantity", "Description", "Unit Cost", "Total Cost", "Date Purchased", "How Purchased") AS
SELECT count(name)                          AS "Quantity",
       name                                 AS "Description",
       cost                                 AS "Unit Cost",
       cost * count(name)::double precision AS "Total Cost",
       date_acquired                        AS "Date Purchased",
       source                               AS "How Purchased"
FROM equipment.equipment
WHERE disposed = false
GROUP BY name, cost, description, date_acquired, source
ORDER BY name;

CREATE VIEW equipment_view
            (name, description, product_num, serial_num, source, cost, date_acquired, hire, disposed, asset_num,
             date_disposed, id, ystv_name, category)
AS
SELECT equipment.name,
       equipment.description,
       equipment.product_num,
       equipment.serial_num,
       equipment.source,
       equipment.cost,
       equipment.date_acquired,
       equipment.hire,
       equipment.disposed,
       equipment.asset_num,
       equipment.date_disposed,
       equipment.id,
       equipment.ystv_name,
       categories.category
FROM equipment.equipment
         JOIN equipment.categories ON equipment.category_id = categories.cat_id;

CREATE VIEW hires_list(category, hire_description, id, item, hire_cost) AS
SELECT hire_categories.category,
       hires.hire_description,
       hires.id,
       hires.item,
       hires.hire_cost
FROM equipment.hire_categories
         JOIN equipment.hires ON hire_categories.cat_id = hires.category_id;

-- +goose Down

DROP SCHEMA equipment CASCADE;