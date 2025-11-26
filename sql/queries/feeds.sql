-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id, created_at, updated_at)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING *;
-- name: GetFeeds :many
SELECT f.id, f.name, f.url, u.name username FROM feeds f, users u WHERE f.user_id = u.id;
-- name: GetFeedByUrl :one
SELECT id, name, url FROM feeds WHERE url = $1 LIMIT 1;
