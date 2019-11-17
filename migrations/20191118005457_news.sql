-- +goose Up

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS news
(
    id         SERIAL PRIMARY KEY,
    title      CITEXT NOT NULL UNIQUE,
    payload    CITEXT,
    created    timestamptz DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE news;
-- +goose StatementEnd
-- SQL in this section is executed when the migration is rolled back.
