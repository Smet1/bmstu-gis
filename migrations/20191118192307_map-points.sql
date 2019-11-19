-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS map_points
(
    id      SERIAL PRIMARY KEY,
    name    citext,
    cabinet bool NOT NULL DEFAULT FALSE,
    stair   bool NOT NULL DEFAULT FALSE,
    x       int  NOT NULL,
    y       int  NOT NULL,
    level   int  NOT NULL,
    node_id int  NOT NULL REFERENCES map_points (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE map_points CASCADE;
-- +goose StatementEnd
