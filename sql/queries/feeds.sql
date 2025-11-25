-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id)
VALUES (
  $1,
  $2,
  $3,
  $4
)
RETURNING *;
-- name: GetFeeds :many
SELECT f.id, f.name, f.url, u.name username FROM feeds f, users u WHERE f.user_id = u.id;
