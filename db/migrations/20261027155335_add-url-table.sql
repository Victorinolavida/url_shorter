-- migrate:up
DO $$
BEGIN
    IF to_regclass('public.links') IS NULL THEN
        CREATE TABLE IF NOT EXISTS url_code (
            id      bigserial PRIMARY KEY,
            code    TEXT NOT NULL UNIQUE,
            url     TEXT NOT NULL,
            created_at    TIMESTAMPTZ not null default now()
        );
    END IF;
END $$;


-- migrate:down
DROP TABLE IF EXISTS url_code;
