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
    category_caption character varying(128) NOT NULL,
    category_status item_status NOT NULL,
    CONSTRAINT catalog_categories_pkey PRIMARY KEY (category_id)
);

-- ALTER SEQUENCE catalog_categories_category_id_seq RESTART WITH 1;


CREATE TABLE catalog_trademarks (
    trademark_id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    datetime_created timestamp NOT NULL DEFAULT now(),
    trademark_caption character varying(64) NOT NULL,
    trademark_status item_status NOT NULL,
    CONSTRAINT catalog_trademarks_pkey PRIMARY KEY (trademark_id)
);

-- ALTER SEQUENCE catalog_categories_trademark_id_seq RESTART WITH 1;


CREATE TABLE catalog_products (
    product_id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    datetime_created timestamp NOT NULL DEFAULT now(),
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
