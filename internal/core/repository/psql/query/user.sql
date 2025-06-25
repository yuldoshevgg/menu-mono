-- name: CreateUser :exec
INSERT INTO users (
    username,
    password,
    role_id,
    venue_id,
    first_name,
    last_name,
    phone_number,
    middle_name
) VALUES (
    sqlc.arg('username'),
    sqlc.arg('password'),
    sqlc.arg('role_id'),
    sqlc.arg('venue_id'),
    sqlc.arg('first_name'),
    sqlc.arg('last_name'),
    sqlc.arg('phone_number'),
    sqlc.arg('middle_name')
);

-- name: GetUser :one
SELECT * FROM users WHERE id = $1 LIMIT 1 AND deleted_at IS NULL;

-- name: UpdateUser :exec
UPDATE users SET 
    username = COALESCE(sqlc.narg('username'), username),
    password = COALESCE(sqlc.narg('password'), password),
    role_id = COALESCE(sqlc.narg('role_id'), role_id),
    venue_id = COALESCE(sqlc.narg('venue_id'), venue_id),
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
    middle_name = COALESCE(sqlc.narg('middle_name'), middle_name) 
WHERE id = sqlc.arg('id') AND deleted_at IS NULL;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1 AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT * FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users WHERE deleted_at IS NULL;