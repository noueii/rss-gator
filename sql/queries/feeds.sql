-- name: CreateFeed :one
INSERT INTO feeds(name, url, user_id)
	VALUES ($1, $2, $3) RETURNING *;

-- name: GetFeedsWithAuthor :many
SELECT feeds.name, feeds.url, users.name
FROM feeds
INNER JOIN users ON users.id = feeds.user_id;

-- name: GetFeedByURL :one
SELECT * FROM feeds
WHERE url = $1
LIMIT 1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
	ORDER BY last_fetched_at ASC NULLS FIRST
	LIMIT 1;
