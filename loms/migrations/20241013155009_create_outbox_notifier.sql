-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS outbox.notifier
(
  id         BIGSERIAL PRIMARY KEY,
  order_id   BIGINT,
  status     outbox.order_status NOT NULL,
  is_sent    BOOLEAN             NOT NULL DEFAULT false,
  created_at TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS outbox.notifier;
-- +goose StatementEnd
