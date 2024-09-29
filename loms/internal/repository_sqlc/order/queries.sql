-- name: CreateOrder :one
INSERT INTO orders.orders (user_id, status, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: GetOrderById :one
SELECT id, user_id, status, created_at, updated_at
FROM orders.orders
WHERE id = $1;

-- name: InsertOrderItem :exec
INSERT INTO orders.order_items (sku_id, order_id, count, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetOrderItems :many
SELECT sku_id, count
FROM orders.order_items
WHERE order_id = $1;

-- name: UpdateOrderStatus :exec
UPDATE orders.orders
SET status = $2, updated_at = $3
WHERE id = $1;
