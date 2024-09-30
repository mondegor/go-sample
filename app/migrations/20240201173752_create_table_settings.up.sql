-- --------------------------------------------------------------------------------------------------

CREATE TABLE sample_catalog.settings (
    setting_id int4 NOT NULL CONSTRAINT pk_settings PRIMARY KEY,
    setting_name character varying(64) NOT NULL,
    setting_type int2 NOT NULL, -- 1=STRING, 2=STRING_LIST, 3=INTEGER, 4=INTEGER_LIST, 5=BOOLEAN
    setting_value character varying(65536) NOT NULL,
    setting_description character varying(1024) NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uk_settings_setting_name ON sample_catalog.settings (setting_name);
CREATE INDEX ix_settings_updated_at ON sample_catalog.settings (updated_at);

-- --------------------------------------------------------------------------------------------------

INSERT INTO sample_catalog.settings (setting_id, setting_name, setting_type, setting_value, setting_description, created_at, updated_at)
VALUES
    (1, 'catalog.categories.list.enabled', 5/*BOOLEAN*/, 'true', 'Отображение списка категорий', '2023-01-01 12:15:59.981966', '2023-01-01 12:15:59.981966'),
    (2, 'catalog.categories.get.enabled', 5/*BOOLEAN*/, 'true', 'Отображение категории', '2023-01-01 12:15:59.981966', '2023-01-01 12:15:59.981966');