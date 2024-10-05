-- name: CreateOrder :one
INSERT INTO orders.orders (user_id, status, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: InsertOrderItem :exec
INSERT INTO orders.order_items (sku_id, order_id, count, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateOrderStatus :exec
UPDATE orders.orders
SET status = $2, updated_at = $3
WHERE id = $1;

-- name: GetOrderWithItems :many
SELECT
  o.id, o.user_id, o.status, o.created_at, o.updated_at,
  CAST(COALESCE(oi.sku_id, 0) AS BIGINT) AS sku_id,
  CAST(COALESCE(oi.count, 0) AS BIGINT) AS count
FROM
  orders.orders o
    LEFT JOIN
  orders.order_items oi ON o.id = oi.order_id
WHERE
  o.id = $1;
