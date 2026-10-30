
-- name: SaveLink :one
INSERT INTO links(
    code,
    raw_url,
    url_hash,
    expires_at
) VALUES ($1,$2,$3,$4)
ON CONFLICT (url_hash) WHERE status = 'active'
DO UPDATE SET url_hash = EXCLUDED.url_hash
RETURNING *;


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
SET status = 'expired',
updated_at = now()
WHERE url_hash = $1
  AND status = 'active'
  AND expires_at IS NOT NULL
  AND expires_at <= now();

-- name: ExpireLinks :many
UPDATE links
SET status = 'expired',
updated_at = now()
WHERE status = 'active'
  AND expires_at IS NOT NULL
  AND expires_at < now()
RETURNING code;

-- name: GetActiveLinkByCode :one
SELECT * from links 
WHERE code = $1 and status = 'active'
  AND (expires_at IS NULL OR expires_at > now());


-- name: CreateClickEvent :one
INSERT INTO click_events (
    link_id,
    user_agent,
    browser,
    os,
    device_type,
    country,
    ip_hash,
    is_bot,
    status_code
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;
