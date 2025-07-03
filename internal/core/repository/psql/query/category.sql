-- name: CreateCategory :exec
INSERT INTO categories (
    name,
    venue_id
) VALUES (
    sqlc.arg('name'),
    sqlc.arg('venue_id')
);

-- name: GetCategory :one
SELECT * FROM categories
WHERE id = sqlc.arg('id')
LIMIT 1;

-- name: ListCategories :many
SELECT * FROM categories
WHERE venue_id = sqlc.arg('venue_id')
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountCategories :one
SELECT COUNT(*) FROM categories
WHERE venue_id = sqlc.arg('venue_id');

-- name: UpdateCategory :exec
UPDATE categories SET 
    name = COALESCE(sqlc.narg('name'), name),
    venue_id = COALESCE(sqlc.narg('venue_id'), venue_id)
WHERE id = sqlc.arg('id')
;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = sqlc.arg('id');