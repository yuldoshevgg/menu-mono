-- name: CreateRole :exec
INSERT INTO roles 
    (title, description) 
VALUES 
    (
        sqlc.arg('title'),
        sqlc.arg('description')
    )
ON CONFLICT (title) DO UPDATE SET 
    description = sqlc.arg('description')
;

-- name: GetRole :one
SELECT * FROM roles WHERE id = $1 LIMIT 1 AND deleted_at IS NULL;

-- name: UpdateRole :exec
UPDATE roles SET 
    title = COALESCE(sqlc.narg('title'), title),
    description = COALESCE(sqlc.narg('description'), description)
WHERE id = sqlc.arg('id') AND deleted_at IS NULL;

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = $1 AND deleted_at IS NULL;

-- name: ListRoles :many
SELECT * FROM roles WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CountRoles :one
SELECT COUNT(*) FROM roles WHERE deleted_at IS NULL;