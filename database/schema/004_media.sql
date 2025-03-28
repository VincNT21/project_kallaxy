-- +goose Up
CREATE TABLE media (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL UNIQUE,
    creator TEXT NOT NULL,
    release_year INTEGER NOT NULL,
    image_url TEXT NOT NULL,
    metadata JSONB NOT NULL
);

-- +goose Down
DROP TABLE media;