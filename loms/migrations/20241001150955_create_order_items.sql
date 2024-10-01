-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS orders.order_items
(
  id         BIGSERIAL PRIMARY KEY,
  sku_id     BIGINT    NOT NULL,
  order_id   BIGINT    NOT NULL REFERENCES orders.orders (id) ON DELETE CASCADE,
  count      BIGINT    NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для ускорения выборок
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON orders.order_items (order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_sku_id ON orders.order_items (sku_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS orders.order_items;
DROP INDEX IF EXISTS idx_order_items_order_id;
DROP INDEX IF EXISTS idx_order_items_sku_id;

-- +goose StatementEnd
