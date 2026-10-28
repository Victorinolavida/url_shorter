-- migrate:up
DO $$
BEGIN
    IF to_regclass('public.url_code') IS NOT NULL
       AND to_regclass('public.links') IS NULL THEN
        ALTER TABLE url_code RENAME TO links;
    END IF;
END $$;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'links' AND column_name = 'url'
    ) THEN
        ALTER TABLE links RENAME COLUMN url TO raw_url;
    END IF;
END $$;

-- migrate:down
DO $$
BEGIN
    IF to_regclass('public.links') IS NOT NULL
       AND to_regclass('public.url_code') IS NULL THEN
        ALTER TABLE links RENAME TO url_code;
    END IF;
END $$;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'url_code' AND column_name = 'raw_url'
    ) THEN
        ALTER TABLE url_code RENAME COLUMN raw_url TO url;
    END IF;
END $$;

