
-- name: SaveLink :one
INSERT INTO links(
    code,
    raw_url,
    url_hash,
    expires_at
) VALUES ($1,$2,$3,$4)
ON CONFLICT (url_hash) WHERE status = 'active'
DO UPDATE SET url_hash = EXCLUDED.url_hash
RETURNING id, code, raw_url, created_at, url_hash, expires_at, status;


-- name: GetLinkByCode :one
SELECT * from links 
WHERE code = $1;


-- name: GetLinkByHash :one
SELECT * from links 
WHERE url_hash = $1;

-- name: GetActiveLinkByHash :one
SELECT * from links 
WHERE url_hash = $1 and status = 'active'
  AND (expires_at IS NULL OR expires_at > now());


-- name: ExpireStaleLinksByHash :exec
UPDATE links
SET status = 'expired'
WHERE url_hash = $1
  AND status = 'active'
  AND expires_at IS NOT NULL
  AND expires_at <= now();

-- name: ExpireLinks :execrows
UPDATE links
SET status = 'expired'
WHERE status = 'active'
  AND expires_at IS NOT NULL
  AND expires_at < now();

-- name: GetActiveLinkByCode :one
SELECT * from links 
WHERE code = $1 and status = 'active'
  AND (expires_at IS NULL OR expires_at > now());
