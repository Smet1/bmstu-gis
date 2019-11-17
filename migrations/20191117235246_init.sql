-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS CITEXT;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,
    login      CITEXT NOT NULL UNIQUE,
    password   TEXT,
    registered timestamptz DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table users;
-- +goose StatementEnd
-- SQL in this section is executed when the migration is rolled back.
