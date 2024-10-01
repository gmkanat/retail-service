// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package pgstocksqry

import (
	"context"
)

const getReservedStock = `-- name: GetReservedStock :one
SELECT reserved
FROM stocks.stocks
WHERE id = $1
`

func (q *Queries) GetReservedStock(ctx context.Context, id int64) (int64, error) {
	row := q.db.QueryRow(ctx, getReservedStock, id)
	var reserved int64
	err := row.Scan(&reserved)
	return reserved, err
}

const getStockBySKU = `-- name: GetStockBySKU :one
SELECT available, reserved
FROM stocks.stocks
WHERE id = $1
`

type GetStockBySKURow struct {
	Available int64
	Reserved  int64
}

func (q *Queries) GetStockBySKU(ctx context.Context, id int64) (GetStockBySKURow, error) {
	row := q.db.QueryRow(ctx, getStockBySKU, id)
	var i GetStockBySKURow
	err := row.Scan(&i.Available, &i.Reserved)
	return i, err
}

const releaseStock = `-- name: ReleaseStock :exec
UPDATE stocks.stocks
SET available = available + $2, reserved = reserved - $2
WHERE id = $1 AND reserved >= $2
`

type ReleaseStockParams struct {
	ID        int64
	Available int64
}

func (q *Queries) ReleaseStock(ctx context.Context, arg ReleaseStockParams) error {
	_, err := q.db.Exec(ctx, releaseStock, arg.ID, arg.Available)
	return err
}

const reserveRemoveStock = `-- name: ReserveRemoveStock :exec
UPDATE stocks.stocks
SET reserved = reserved - $2
WHERE id = $1 AND reserved >= $2
`

type ReserveRemoveStockParams struct {
	ID       int64
	Reserved int64
}

func (q *Queries) ReserveRemoveStock(ctx context.Context, arg ReserveRemoveStockParams) error {
	_, err := q.db.Exec(ctx, reserveRemoveStock, arg.ID, arg.Reserved)
	return err
}

const reserveStock = `-- name: ReserveStock :exec
UPDATE stocks.stocks
SET available = available - $2, reserved = reserved + $2
WHERE id = $1 AND available >= $2
`

type ReserveStockParams struct {
	ID        int64
	Available int64
}

func (q *Queries) ReserveStock(ctx context.Context, arg ReserveStockParams) error {
	_, err := q.db.Exec(ctx, reserveStock, arg.ID, arg.Available)
	return err
}
