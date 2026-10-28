-- migrate:up
-- status values active, expired, disabled, blocked
ALTER TABLE links ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'active';
UPDATE links SET status = 'active' WHERE status IS NULL;
ALTER TABLE links ALTER COLUMN status SET NOT NULL;
ALTER TABLE links ALTER COLUMN status SET DEFAULT 'active';

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chck_status' AND conrelid = 'links'::regclass
    ) THEN
        ALTER TABLE links ADD CONSTRAINT chck_status
            CHECK (status IN ('active', 'expired', 'disabled', 'blocked'));
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'unique_link_code' AND conrelid = 'links'::regclass
    ) THEN
        ALTER TABLE links ADD CONSTRAINT unique_link_code UNIQUE(code);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chck_raw_url_length' AND conrelid = 'links'::regclass
    ) THEN
        ALTER TABLE links ADD CONSTRAINT chck_raw_url_length CHECK (length(raw_url) <= 2048);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chck_code_length' AND conrelid = 'links'::regclass
    ) THEN
        ALTER TABLE links ADD CONSTRAINT chck_code_length CHECK ((char_length(code) BETWEEN 4 AND 32));
    END IF;
END $$;

-- migrate:down
DO $$
BEGIN
    IF to_regclass('public.links') IS NOT NULL THEN
        ALTER TABLE links DROP COLUMN IF EXISTS status;
        ALTER TABLE links DROP CONSTRAINT IF EXISTS unique_link_code;
        ALTER TABLE links DROP CONSTRAINT IF EXISTS chck_raw_url_length;
        ALTER TABLE links DROP CONSTRAINT IF EXISTS chck_code_length;
        ALTER TABLE links DROP CONSTRAINT IF EXISTS chck_status;
    END IF;
END $$;
