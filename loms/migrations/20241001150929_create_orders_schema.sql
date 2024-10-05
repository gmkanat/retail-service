-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS orders;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS orders;
-- +goose StatementEnd
