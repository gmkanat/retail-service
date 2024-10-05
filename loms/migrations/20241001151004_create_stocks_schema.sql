-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA IF NOT EXISTS stocks;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP SCHEMA IF EXISTS stocks;

-- +goose StatementEnd
