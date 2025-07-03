-- name: CreateMenuItem :exec
INSERT INTO menu_items (
    title,
    category_id,
    venue_id,
    price,
    image_url,
    is_available,
    description
) VALUES (
    sqlc.arg('title'),
    sqlc.arg('category_id'),
    sqlc.arg('venue_id'),
    sqlc.arg('price'),
    sqlc.arg('image_url'),
    sqlc.arg('is_available'),
    sqlc.arg('description')
);

-- name: GetMenuItem :one
SELECT * FROM menu_items
WHERE id = sqlc.arg('id')
LIMIT 1;

-- name: ListMenuItems :many
SELECT * FROM menu_items
WHERE venue_id = sqlc.arg('venue_id')
AND category_id = sqlc.arg('category_id')
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountMenuItems :one
SELECT COUNT(*) FROM menu_items
WHERE venue_id = sqlc.arg('venue_id');

-- name: UpdateMenuItem :exec
UPDATE menu_items SET 
    title = COALESCE(sqlc.narg('title'), title),
    category_id = COALESCE(sqlc.narg('category_id'), category_id),
    venue_id = COALESCE(sqlc.narg('venue_id'), venue_id),
    price = COALESCE(sqlc.narg('price'), price),
    image_url = COALESCE(sqlc.narg('image_url'), image_url),
    is_available = COALESCE(sqlc.narg('is_available'), is_available),
    description = COALESCE(sqlc.narg('description'), description)
WHERE id = sqlc.arg('id')
;

-- name: DeleteMenuItem :exec
DELETE FROM menu_items WHERE id = sqlc.arg('id');

-- name: CreateQrLink :one
INSERT INTO qr_links (
    venue_id,
    table_number,
    expires_at
) VALUES (
    sqlc.arg('venue_id'),
    sqlc.arg('table_number'),
    sqlc.arg('expires_at')
) RETURNING id;