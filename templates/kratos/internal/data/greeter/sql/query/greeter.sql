-- name: CreateGreeter :execresult
INSERT INTO greeters (hello) VALUES (?);

-- name: UpdateGreeter :exec
UPDATE greeters
SET hello = ?
WHERE id = ?;

-- name: GetGreeterByID :one
SELECT id, hello, created_at
FROM greeters
WHERE id = ?
LIMIT 1;

-- name: ListGreetersByHello :many
SELECT id, hello, created_at
FROM greeters
WHERE hello = ?
ORDER BY id DESC;

-- name: ListGreeters :many
SELECT id, hello, created_at
FROM greeters
ORDER BY id DESC;
