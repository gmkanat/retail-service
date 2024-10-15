-- +goose Up
-- +goose StatementBegin
CREATE TYPE outbox.order_status AS ENUM (
  'New',
  'AwaitingPayment',
  'Paid',
  'Cancelled',
  'Failed'
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS outbox.order_status;
-- +goose StatementEnd
