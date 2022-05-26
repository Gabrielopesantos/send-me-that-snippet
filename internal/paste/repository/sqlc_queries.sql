-- name: Insert :one
INSERT INTO pastes (
    id, content, content_sha, language, created_at, expires_in, expired
) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: Get :one
SELECT * FROM pastes
WHERE id = $1;

-- name: Delete :one
DELETE FROM pastes
WHERE id = $1
RETURNING *;

-- name: Update :one
UPDATE pastes
SET
    content = $1,
    content_sha = $2,
    language = $3,
    expired = $4
WHERE id = $5
RETURNING *;

-- -- name: Update :one
-- UPDATE pastes
-- SET
--     content = COALESCE(NULLIF($1,''), content),
--     content_sha = COALESCE(NULLIF($2, ''), content_sha),
--     language = COALESCE(NULLIF($3, ''), language),
--     expired = COALESCE(NULLIF($4, false), expired)
-- WHERE id = $5
-- RETURNING *;
