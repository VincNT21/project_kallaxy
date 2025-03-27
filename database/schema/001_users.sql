-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    username TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL DEFAULT('unset'),
    email TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;