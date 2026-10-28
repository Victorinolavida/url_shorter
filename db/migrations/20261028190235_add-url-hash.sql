-- migrate:up
DO $$
BEGIN
    IF to_regclass('public.url_code') IS NOT NULL THEN
        ALTER TABLE url_code ADD COLUMN IF NOT EXISTS url_hash bytea not null;
        ALTER TABLE url_code ADD COLUMN IF NOT EXISTS expires_at TIMESTAMPTZ default null;
    END IF;
END $$;

-- migrate:down

DO $$
BEGIN
    IF to_regclass('public.url_code') IS NOT NULL THEN
        ALTER TABLE url_code DROP COLUMN IF EXISTS url_hash;
        ALTER TABLE url_code DROP COLUMN IF EXISTS expires_at;
    END IF;
END $$;
