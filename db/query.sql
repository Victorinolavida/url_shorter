
-- name: SaveLink :one
INSERT INTO links(
    code,
    raw_url,
    url_hash,
    expires_at
) VALUES ($1,$2,$3,$4)
RETURNING  id, created_at;


-- name: GetLinkByCode :one
SELECT * from links 
WHERE code = $1;


-- name: GetLinkByHash :one
SELECT * from links 
WHERE url_hash = $1;

-- name: GetActiveLinkByHash :one
SELECT * from links 
WHERE url_hash = $1 and status = 'active';
