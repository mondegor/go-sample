-- --------------------------------------------------------------------------------------------------

CREATE TABLE sample_catalog.categories (
    category_id uuid NOT NULL CONSTRAINT pk_categories PRIMARY KEY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    category_caption character varying(128) NOT NULL,
    image_meta jsonb DEFAULT NULL,
    category_status int2 NOT NULL, -- 1=DRAFT, 2=ENABLED, 3=DISABLED
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NULL,
    deleted_at timestamp with time zone NULL
);

-- --------------------------------------------------------------------------------------------------

INSERT INTO sample_catalog.categories (category_id, tag_version, category_caption, image_meta, category_status, created_at, updated_at, deleted_at)
VALUES
    ('b86555ab-9320-4680-b62d-1ea449550fff', 1, 'Электроника', NULL, 2/*ENABLED*/, '2023-01-01 12:15:59.981966', NULL, NULL),
    ('166a72b5-b9fa-499c-8140-3627b34acbbe', 1, 'Бытовая техника', NULL, 2/*ENABLED*/, '2023-02-01 13:15:08.678538', NULL, NULL),
    ('3979b1bb-26ba-4a57-b9ee-fd723b4fa9a0', 1, 'Компьютерная техника', NULL, 1/*DRAFT*/, '2023-03-01 14:15:16.078388', NULL, NULL),
    ('75b3e9ea-7c07-4bbb-a8b1-adb02ebcbbec', 1, 'Строительство и ремонт', NULL, 2/*ENABLED*/, '2023-04-01 15:15:12.319733', NULL, NULL),
    ('0a2b96df-a528-4081-8c76-3280ffba320e', 1, 'Кондитерские изделия', NULL, 2/*ENABLED*/, '2023-05-01 16:15:14.319733', NULL, NULL),
    ('d1c7ddd8-6d48-4008-be0e-d4170ef62f1c', 1, 'Увлажнители воздуха', NULL, 3/*DISABLED*/, '2023-06-01 17:15:32.319733', NULL, NULL),
    ('fb0ebfa1-d264-46d0-a377-c3b8f67b01a3', 1, 'Тостеры и ростеры', NULL, 1/*DRAFT*/, '2023-07-01 18:15:23.319733', NULL, NULL),
    ('af6f7154-6ecc-4734-983b-1b91a866a6bc', 1, 'Сварочные аппараты', NULL, 2/*ENABLED*/, '2023-08-01 19:15:19.319733', NULL, NULL);