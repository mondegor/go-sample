-- --------------------------------------------------------------------------------------------------

CREATE TABLE sample_catalog.products (
    product_id int8 NOT NULL GENERATED BY DEFAULT AS IDENTITY CONSTRAINT pk_products PRIMARY KEY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    category_id uuid NOT NULL CONSTRAINT fk_products_category_id REFERENCES sample_catalog.categories (category_id),
    product_article character varying(32) NULL,
    product_caption character varying(128) NOT NULL,
    trademark_id int8 NOT NULL CONSTRAINT fk_products_trademark_id REFERENCES sample_catalog.trademarks (trademark_id),
    product_price int8 NOT NULL, -- coins * 100
    product_status int2 NOT NULL, -- 1=DRAFT, 2=ENABLED, 3=DISABLED
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
    deleted_at timestamp with time zone NULL,
    prev_field_id int8 NULL CHECK(prev_field_id IS NULL OR prev_field_id > 0),
    next_field_id int8 NULL CHECK(next_field_id IS NULL OR next_field_id > 0),
    order_index int8 NULL CHECK(order_index IS NULL OR order_index > 0)
);

CREATE UNIQUE INDEX uk_products_product_article ON sample_catalog.products (product_article) WHERE deleted_at IS NULL;

-- --------------------------------------------------------------------------------------------------

INSERT INTO sample_catalog.products (product_id, tag_version, category_id, product_article, product_caption, trademark_id, product_price, product_status, created_at, updated_at, deleted_at, prev_field_id, next_field_id, order_index)
VALUES
    (1, 1, '0a2b96df-a528-4081-8c76-3280ffba320e', 'candy-001', 'Алёнка', 1, 9050, 2/*ENABLED*/, '2023-01-01 12:30:59.981966', '2023-01-01 12:30:59.981966', NULL, NULL, 3, 1000),
    (2, 1, 'af6f7154-6ecc-4734-983b-1b91a866a6bc', 'masha-002', 'Маша и Медведь', 2, 5380000, 2/*ENABLED*/, '2023-02-01 13:30:08.678538', '2023-02-01 13:30:08.678538', NULL, NULL, NULL, NULL),
    (3, 1, '0a2b96df-a528-4081-8c76-3280ffba320e', 'cookie-003', 'Юбилейное', 3, 4300, 2/*ENABLED*/, '2023-03-01 14:30:16.078388', '2023-03-01 14:30:16.078388', NULL, 1, 4, 2000),
    (4, 1, '0a2b96df-a528-4081-8c76-3280ffba320e', 'candy-004', 'Бабаевский', 4, 12570, 2/*ENABLED*/, '2023-04-01 15:30:23.319733', '2023-04-01 15:30:23.319733', NULL, 3, 5, 3000),
    (5, 1, '0a2b96df-a528-4081-8c76-3280ffba320e', 'candy-005', 'Кара-Кум', 5, 5000, 1/*DRAFT*/, '2023-05-01 16:30:23.319733', '2023-05-01 16:30:23.319733', NULL, 4, 6, 4000),
    (6, 1, '0a2b96df-a528-4081-8c76-3280ffba320e', 'candy-006', 'Мишка на Севере', 6, 21450, 2/*ENABLED*/, '2023-06-01 17:30:14.319733', '2023-06-01 17:30:14.319733', NULL, 5, NULL, 5000),
    (7, 1, 'b86555ab-9320-4680-b62d-1ea449550fff', 'tv-007', 'Радуга ТВ', 7, 2000000, 2/*ENABLED*/, '2023-07-01 18:30:23.319733', '2023-07-01 18:30:23.319733', NULL, NULL, 8, 1000),
    (8, 1, 'b86555ab-9320-4680-b62d-1ea449550fff', 'rocket-008', 'Ракета', 8, 1530000, 1/*DRAFT*/, '2023-08-01 19:30:45.319733', '2023-08-01 19:30:45.319733', NULL, 8, NULL, 2000);

CREATE INDEX ix_products_category_id_order_index ON sample_catalog.products (category_id, order_index) WHERE deleted_at IS NULL;

ALTER SEQUENCE sample_catalog.products_product_id_seq RESTART WITH 9;