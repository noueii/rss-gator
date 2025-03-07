-- name: CreateFeedFollow :one
INSERT INTO feed_follows (user_id, feed_id)
VALUES ($1, $2) 
RETURNING *,
	(SELECT name FROM users WHERE users.id = feed_follows.user_id) as user_name,
	(SELECT name FROM feeds WHERE feeds.id = feed_follows.feed_id) as feed_name;

-- name: GetFeedFollowsForUser :many
SELECT *, feeds.name AS feed_name FROM feed_follows
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;

-- name: DeleteUserFeed :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;

