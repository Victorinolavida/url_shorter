-- migrate:up
ALTER TABLE links ADD COLUMN updated_at Timestamptz DEFAULT NULL;

-- migrate:down
ALTER TABLE links DROP COLUMN updated_at;

