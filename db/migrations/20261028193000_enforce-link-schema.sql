-- migrate:up
-- Make id a GENERATED ALWAYS AS IDENTITY primary key, preserving existing ids.
DO $$
DECLARE
    seq_name text;
    next_id bigint;
    is_identity boolean;
BEGIN
    SELECT attidentity <> '' INTO is_identity
    FROM pg_attribute
    WHERE attrelid = 'links'::regclass AND attname = 'id' AND NOT attisdropped;

    IF NOT COALESCE(is_identity, false) THEN
        SELECT COALESCE(MAX(id), 0) + 1 INTO next_id FROM links;
        seq_name := pg_get_serial_sequence('links', 'id');

        ALTER TABLE links ALTER COLUMN id DROP DEFAULT;

        IF seq_name IS NOT NULL THEN
            EXECUTE format('ALTER SEQUENCE %s OWNED BY NONE', seq_name);
            EXECUTE format('DROP SEQUENCE %s', seq_name);
        END IF;

        ALTER TABLE links ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY;
        EXECUTE format('ALTER TABLE links ALTER COLUMN id RESTART WITH %s', next_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'url_code_pkey' AND conrelid = 'links'::regclass
    ) THEN
        ALTER TABLE links ADD CONSTRAINT url_code_pkey PRIMARY KEY (id);
    END IF;
END $$;

-- Ensure a single UNIQUE (code) constraint so ON CONFLICT (code) works.
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'unique_link_code' AND conrelid = 'links'::regclass
    ) THEN
        ALTER TABLE links ADD CONSTRAINT unique_link_code UNIQUE (code);
    END IF;

    IF EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'url_code_code_key' AND conrelid = 'links'::regclass
    ) THEN
        ALTER TABLE links DROP CONSTRAINT url_code_code_key;
    END IF;
END $$;

-- status must never be NULL and only accept known values.
UPDATE links SET status = 'active' WHERE status IS NULL;
ALTER TABLE links ALTER COLUMN status SET DEFAULT 'active';
ALTER TABLE links ALTER COLUMN status SET NOT NULL;

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


-- migrate:down
DO $$
DECLARE
    next_id bigint;
    is_identity boolean;
BEGIN
    IF to_regclass('public.links') IS NULL THEN
        RETURN;
    END IF;

    ALTER TABLE links DROP CONSTRAINT IF EXISTS chck_status;
    ALTER TABLE links ALTER COLUMN status DROP NOT NULL;

    SELECT attidentity <> '' INTO is_identity
    FROM pg_attribute
    WHERE attrelid = 'links'::regclass AND attname = 'id' AND NOT attisdropped;

    IF COALESCE(is_identity, false) THEN
        SELECT COALESCE(MAX(id), 0) + 1 INTO next_id FROM links;
        ALTER TABLE links ALTER COLUMN id DROP IDENTITY IF EXISTS;
        CREATE SEQUENCE IF NOT EXISTS url_code_id_seq OWNED BY links.id;
        PERFORM setval('url_code_id_seq', next_id, false);
        ALTER TABLE links ALTER COLUMN id
            SET DEFAULT nextval('url_code_id_seq'::regclass);
    END IF;
END $$;
