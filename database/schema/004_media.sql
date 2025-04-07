-- +goose Up
CREATE TABLE media (
    id UUID PRIMARY KEY,
    media_type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    creator TEXT NOT NULL,
    release_year TEXT NOT NULL,
    image_url TEXT NOT NULL,
    metadata JSONB NOT NULL,
    UNIQUE (media_type, title)
);

-- +goose Down
DROP TABLE media;