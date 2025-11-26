-- name: FollowFeed :one
WITH inserted_feed_follow AS (
INSERT INTO feed_follows (id, feed_id, user_id, created_at, updated_at)
VALUES($1,$2,$3,$4,$5)
RETURNING id, created_at, updated_at, user_id, feed_id
)
SELECT inserted_feed_follow.id, inserted_feed_follow.created_at, inserted_feed_follow.updated_at, inserted_feed_follow.user_id, inserted_feed_follow.feed_id,
feeds.name feed_name,
users.name user_name
FROM inserted_feed_follow
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users ON inserted_feed_follow.user_id = users.id;
-- name: GetFeedFollowsForUser :many
SELECT feeds.id, feeds.name, feeds.url, users.name
FROM feed_follows
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
INNER JOIN users ON feed_follows.user_id = users.id
WHERE feed_follows.user_id = $1;
