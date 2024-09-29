-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA IF NOT EXISTS orders;

CREATE TYPE orders.order_status AS ENUM (
  'New',
  'AwaitingPayment',
  'Paid',
  'Cancelled',
  'Failed'
  );

CREATE TABLE IF NOT EXISTS orders.orders
(
  id         BIGSERIAL PRIMARY KEY,
  user_id    BIGINT              NOT NULL,
  status     orders.order_status NOT NULL DEFAULT 'New',
  created_at TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orders.order_items
(
  sku_id     BIGINT    NOT NULL,
  order_id   BIGINT    NOT NULL REFERENCES orders.orders (id) ON DELETE CASCADE,
  count      BIGINT    NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
  );


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS orders.order_items;
DROP TABLE IF EXISTS orders.orders;

DROP TYPE IF EXISTS orders.order_status;
DROP SCHEMA IF EXISTS orders;

-- +goose StatementEnd
