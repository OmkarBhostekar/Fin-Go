-- name: CreateEntry :one
INSERT INTO entries(
    account_id,
    amount,
) VALUES (
    $1,
    $2
) RETURNING *;

-- name: GetEntryById :one
SELECT * FROM entries WHERE id = $1 LIMIT 1;

-- name: GetAllEntries :many
SELECT * FROM entries
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateEntryAmount :exec
UPDATE entries
SET amount = $1
WHERE id = $2

-- name: DeleteEntryById :exec
DELETE FROM entries
WHERE id = $1;
