-- name: CreateVenue :exec
INSERT INTO venues (title, slug, venue_id, logo_url, status, description) VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetVenue :one
SELECT * FROM venues WHERE id = $1 LIMIT 1 AND deleted_at IS NULL;

-- name: UpdateVenue :exec
UPDATE venues SET 
    title = COALESCE(sqlc.narg('title'), title),
    slug = COALESCE(sqlc.narg('slug'), slug),
    venue_id = COALESCE(sqlc.narg('venue_id'), venue_id),
    logo_url = COALESCE(sqlc.narg('logo_url'), logo_url),
    status = COALESCE(sqlc.narg('status'), status),
    description = COALESCE(sqlc.narg('description'), description)
WHERE id = sqlc.arg('id') AND deleted_at IS NULL;

-- name: DeleteVenue :exec
DELETE FROM venues WHERE id = $1 AND deleted_at IS NULL;

-- name: ListVenues :many
SELECT * FROM venues WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CountVenues :one
SELECT COUNT(*) FROM venues WHERE deleted_at IS NULL;

-- name: ExistVenue :one
SELECT EXISTS(SELECT 1 FROM venues WHERE id = $1 AND deleted_at IS NULL);