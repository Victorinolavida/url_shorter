
-- name: SaveUrlCode :one
INSERT INTO url_code(
    code,
    url
) VALUES ($1,$2)
RETURNING  id, created_at;


-- name: GetByCode :one
SELECT * from url_code 
WHERE code = $1;


-- name: GetByURL :one
SELECT * from url_code 
WHERE url = $1;
