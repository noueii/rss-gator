-- name: CreatePost :one
INSERT INTO posts(title, url, description, published_at, feed_id) 
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetPostsForUserWithLimit :many
SELECT p.* FROM feed_follows ff
INNER JOIN posts p ON ff.feed_id = p.feed_id
WHERE ff.user_id = $1
ORDER BY p.published_at
LIMIT $2;

