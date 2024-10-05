-- +goose Up
-- +goose StatementBegin

CREATE TYPE orders.order_status AS ENUM (
  'New',
  'AwaitingPayment',
  'Paid',
  'Cancelled',
  'Failed'
  );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TYPE IF EXISTS orders.order_status;

-- +goose StatementEnd
