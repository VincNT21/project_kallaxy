-- +goose Up
CREATE TABLE users_media_records (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    media_id UUID REFERENCES media(id) ON DELETE CASCADE,
    is_finished BOOLEAN,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    duration INTERVAL DAY,
    UNIQUE (user_id, media_id)
);

-- +goose Down
DROP TABLE users_media_records;