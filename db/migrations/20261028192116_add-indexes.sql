-- migrate:up
CREATE INDEX IF NOT EXISTS idx_link_expired_at
ON links (expires_at)
WHERE expires_at IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_url_hash_expired_at
ON links (url_hash);


-- migrate:down
DROP INDEX IF EXISTS idx_link_expired_at;
DROP INDEX IF EXISTS idx_url_hash_expired_at;
