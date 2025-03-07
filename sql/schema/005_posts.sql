-- +goose Up
CREATE TABLE posts (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	url TEXT UNIQUE NOT NULL,
	published_at TIMESTAMP NOT NULL,
	feed_id UUID NOT NULL,
	
	CONSTRAINT fk_feed FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE 
);

-- +goose Down
DROP TABLE posts;
