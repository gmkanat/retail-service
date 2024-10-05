-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS orders.orders
(
  id         BIGSERIAL PRIMARY KEY,
  user_id    BIGINT              NOT NULL,
  status     orders.order_status NOT NULL DEFAULT 'New',
  created_at TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS orders.orders;

-- +goose StatementEnd
