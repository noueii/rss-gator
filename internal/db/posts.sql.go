// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: posts.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts(title, url, description, published_at, feed_id) 
VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at, title, description, url, published_at, feed_id
`

type CreatePostParams struct {
	Title       string
	Url         string
	Description string
	PublishedAt time.Time
	FeedID      uuid.UUID
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.PublishedAt,
		arg.FeedID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Description,
		&i.Url,
		&i.PublishedAt,
		&i.FeedID,
	)
	return i, err
}

const getPostsForUserWithLimit = `-- name: GetPostsForUserWithLimit :many
SELECT p.id, p.created_at, p.updated_at, p.title, p.description, p.url, p.published_at, p.feed_id FROM feed_follows ff
INNER JOIN posts p ON ff.feed_id = p.feed_id
WHERE ff.user_id = $1
ORDER BY p.published_at
LIMIT $2
`

type GetPostsForUserWithLimitParams struct {
	UserID uuid.UUID
	Limit  int32
}

func (q *Queries) GetPostsForUserWithLimit(ctx context.Context, arg GetPostsForUserWithLimitParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsForUserWithLimit, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Description,
			&i.Url,
			&i.PublishedAt,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
