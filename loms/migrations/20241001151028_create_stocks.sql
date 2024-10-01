-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS stocks.stocks
(
  id         BIGINT PRIMARY KEY,
  available  BIGINT    NOT NULL,
  reserved   BIGINT    NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS stocks.stocks;

-- +goose StatementEnd
