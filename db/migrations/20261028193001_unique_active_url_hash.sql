-- migrate:up
-- Keep only the oldest active link per url_hash before enforcing uniqueness.
WITH duplicates AS (
    SELECT id, row_number() OVER (PARTITION BY url_hash ORDER BY created_at, id) AS rn
    FROM links
    WHERE status = 'active'
)
UPDATE links l
SET status = 'expired'
FROM duplicates d
WHERE l.id = d.id
  AND d.rn > 1;

CREATE UNIQUE INDEX IF NOT EXISTS unique_active_url_hash
ON links (url_hash)
WHERE status = 'active';

-- migrate:down
DROP INDEX IF EXISTS unique_active_url_hash;
