-- +goose Up
-- +goose StatementBegin
ALTER TABLE metrics ALTER COLUMN delta TYPE BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE metrics ALTER COLUMN delta TYPE INTEGER;
-- +goose StatementEnd
