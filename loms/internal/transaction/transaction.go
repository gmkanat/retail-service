package transaction

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/pgcluster"
)

type TxManager struct {
	cluster *pgcluster.Cluster
}

func NewTxManager(cluster *pgcluster.Cluster) *TxManager {
	return &TxManager{cluster: cluster}
}

type txKeyType struct{}

var TxKey = txKeyType{}

func WithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

func GetTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	return tx, ok
}

func (tm *TxManager) WithRepeatableReadTx(ctx context.Context, fn func(ctx context.Context) error) error {
	conn, err := tm.cluster.GetWriter(ctx)
	if err != nil {
		return fmt.Errorf("failed to get writer: %w", err)
	}

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	ctxWithTx := WithTx(ctx, tx)

	if err := fn(ctxWithTx); err != nil {
		return fmt.Errorf("transaction function failed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
