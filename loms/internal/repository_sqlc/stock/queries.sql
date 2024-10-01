-- name: GetStockBySKU :one
SELECT available, reserved
FROM stocks.stocks
WHERE id = $1;

-- name: ReserveStock :exec
UPDATE stocks.stocks
SET available = available - $2, reserved = reserved + $2
WHERE id = $1 AND available >= $2;

-- name: ReleaseStock :exec
UPDATE stocks.stocks
SET available = available + $2, reserved = reserved - $2
WHERE id = $1 AND reserved >= $2;

-- name: ReserveRemoveStock :exec
UPDATE stocks.stocks
SET reserved = reserved - $2
WHERE id = $1 AND reserved >= $2;

-- name: GetReservedStock :one
SELECT reserved
FROM stocks.stocks
WHERE id = $1;
