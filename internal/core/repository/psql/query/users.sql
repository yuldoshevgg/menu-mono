-- name: CreateUser :exec
INSERT INTO users (
    username,
    password,
    role_id,
    restaurant_id,
    full_name
) VALUES (
    sqlc.arg('username'),
    sqlc.arg('password'),
    sqlc.arg('role_id'),
    sqlc.arg('restaurant_id'),
    sqlc.arg('full_name')
);