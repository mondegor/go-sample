-- CREATE SCHEMA public AUTHORIZATION user_pg;

CREATE TYPE item_status AS ENUM (
    'DRAFT',
    'ENABLED',
    'DISABLED',
    'REMOVED');

CREATE TABLE catalog_categories (
    category_id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    datetime_created timestamp NOT NULL DEFAULT now(),
    datetime_updated timestamp NULL,
    category_caption character varying(128) NOT NULL,
    image_path character varying(128) NOT NULL DEFAULT '',
    category_status item_status NOT NULL,
    CONSTRAINT catalog_categories_pkey PRIMARY KEY (category_id)
);

INSERT INTO catalog_categories (category_id, tag_version, datetime_created, datetime_updated, category_caption, image_path, category_status)
    OVERRIDING SYSTEM VALUE
VALUES
    (1, 1, '2023-09-14 15:51:59.981966', NULL, 'Электроника', '', 'DRAFT'),
    (2, 1, '2023-09-14 15:52:08.678538', NULL, 'Бытовая техника', '', 'DRAFT'),
    (3, 1, '2023-09-14 15:52:16.078388', NULL, 'Компьютерная техника', '', 'DRAFT'),
    (4, 1, '2023-09-14 15:52:23.319733', NULL, 'Строительство и ремонт', '', 'DRAFT');

ALTER SEQUENCE catalog_categories_category_id_seq RESTART WITH 5;


CREATE TABLE catalog_trademarks (
    trademark_id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    datetime_created timestamp NOT NULL DEFAULT now(),
    datetime_updated timestamp NULL,
    trademark_caption character varying(64) NOT NULL,
    trademark_status item_status NOT NULL,
    CONSTRAINT catalog_trademarks_pkey PRIMARY KEY (trademark_id)
);

-- ALTER SEQUENCE catalog_trademarks_trademark_id_seq RESTART WITH 1;


CREATE TABLE catalog_products (
    product_id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    datetime_created timestamp NOT NULL DEFAULT now(),
    datetime_updated timestamp NULL,
    category_id int4 NOT NULL,
    trademark_id int4 NOT NULL,
    product_article character varying(32) NULL,
    product_caption character varying(128) NOT NULL,
    product_price int8 NOT NULL CHECK(product_price > 0 AND product_price < 100000000001), -- coins * 100
    product_status item_status NOT NULL,
    prev_field_id int4 NULL CHECK(prev_field_id IS NULL OR prev_field_id > 0),
    next_field_id int4 NULL CHECK(next_field_id IS NULL OR next_field_id > 0),
    order_field int8 NULL CHECK(order_field IS NULL OR order_field > 0),
    CONSTRAINT catalog_products_pkey PRIMARY KEY (product_id),
    CONSTRAINT catalog_products_category_id FOREIGN KEY (category_id)
        REFERENCES catalog_categories(category_id),
    CONSTRAINT catalog_products_trademark_id FOREIGN KEY (trademark_id)
        REFERENCES catalog_trademarks(trademark_id),
    CONSTRAINT catalog_products_product_article UNIQUE (product_article)
);

CREATE INDEX catalog_products_order_field ON catalog_products (category_id, order_field);

-- ALTER SEQUENCE catalog_products_product_id_seq RESTART WITH 1;
