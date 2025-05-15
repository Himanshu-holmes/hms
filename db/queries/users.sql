-- name: CreateUser :one
INSERT INTO users (
    username, password_hash, role, first_name, last_name, email, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 AND is_active = TRUE LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET
    first_name = COALESCE(sqlc.narg(first_name), first_name),
    last_name = COALESCE(sqlc.narg(last_name), last_name),
    email = COALESCE(sqlc.narg(email), email),
    is_active = COALESCE(sqlc.narg(is_active), is_active),
    role = COALESCE(sqlc.narg(role), role),
    password_hash = COALESCE(sqlc.narg(password_hash), password_hash), -- Be careful updating password
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteUser :exec
-- For actual deletion. Consider deactivating (is_active = false) instead.
DELETE FROM users
WHERE id = $1;

-- name: SetUserActiveStatus :one
UPDATE users
SET is_active = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;