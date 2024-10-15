-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS outbox;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS outbox;
-- +goose StatementEnd
