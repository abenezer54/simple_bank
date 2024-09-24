-- name: CreateEntry :one
INSERT INTO entries (account_id, amount)
VALUES ($1, $2)
RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1;

-- name: ListEntriesByAccountID :many
SELECT * FROM entries
WHERE account_id = $1;

-- name: UpdateEntry :one
UPDATE entries
SET account_id = $2, amount = $3, created_at = $4
WHERE id = $1
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1
RETURNING *;