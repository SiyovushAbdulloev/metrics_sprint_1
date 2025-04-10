-- +goose Up
-- +goose StatementBegin
CREATE table IF NOT EXISTS metrics (
    id varchar(255) PRIMARY KEY NOT NULL,
    type varchar(255) NOT NULL,
    delta integer,
    value double precision
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table IF EXISTS metrics;
-- +goose StatementEnd
