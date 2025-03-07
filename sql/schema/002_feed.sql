-- +goose Up

CREATE TABLE feeds(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	name TEXT NOT NULL,
	url TEXT NOT NULL,
	user_id UUID NOT NULL,

	CONSTRAINT fk_user
	FOREIGN KEY (user_id)
	REFERENCES users(id)
	ON DELETE CASCADE,

	CONSTRAINT userId_url UNIQUE (user_id, url)
);

-- +goose Down
DROP TABLE feeds;
