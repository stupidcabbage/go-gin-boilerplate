-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    email varchar(64) primary key,
    username varchar(32),
    "password" text,
    created_at timestamp,
    updated_at timestamp,
    UNIQUE(username)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE users;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
