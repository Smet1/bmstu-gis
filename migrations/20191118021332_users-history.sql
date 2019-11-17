-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users_history
(
    id         SERIAL PRIMARY KEY,
    user_id    int REFERENCES users (id),
    point_from TEXT,
    point_to   TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users_history;
-- +goose StatementEnd
