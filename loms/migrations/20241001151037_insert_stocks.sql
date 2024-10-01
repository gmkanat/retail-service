-- +goose Up
-- +goose StatementBegin

INSERT INTO stocks.stocks (id, available, reserved)
VALUES
  (773297411, 150, 10),
  (1076963, 160, 12),
  (1625903, 170, 14),
  (2956315, 180, 16),
  (1002, 200, 20),
  (1003, 250, 30),
  (1004, 300, 40),
  (1005, 350, 50);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM stocks.stocks
WHERE id IN (773297411, 1076963, 1625903, 2956315, 1002, 1003, 1004, 1005);

-- +goose StatementEnd
