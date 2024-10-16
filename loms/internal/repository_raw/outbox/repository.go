package outbox

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/infra"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/pgcluster"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/transaction"
)

type Repository struct {
	cluster       *pgcluster.Cluster
	kafkaProducer *infra.KafkaProducer
}

func NewRepository(cluster *pgcluster.Cluster, kafkaProducer *infra.KafkaProducer) *Repository {
	return &Repository{
		cluster:       cluster,
		kafkaProducer: kafkaProducer,
	}
}

func (n *Repository) CreateEvent(ctx context.Context, event model.Event) error {
	tx, ok := transaction.GetTx(ctx)
	if !ok {
		return fmt.Errorf("transaction not found in context")
	}

	qry := `INSERT INTO outbox.notifier (order_id, status)
			VALUES ($1, $2);`

	if _, err := tx.Exec(ctx, qry, event.OrderID, event.Status); err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

func (n *Repository) FetchAndMarkBatch(ctx context.Context, batchSize int) ([]*model.Event, error) {
	tx, ok := transaction.GetTx(ctx)
	if !ok {
		return nil, fmt.Errorf("transaction not found in context")
	}

	qry := `
        UPDATE outbox.notifier
        SET is_sent = true
        WHERE id IN (
            SELECT id
            FROM outbox.notifier
            WHERE is_sent = false
            ORDER BY id
            LIMIT $1
            FOR UPDATE SKIP LOCKED
        )
        RETURNING id, order_id, status, created_at`

	rows, err := tx.Query(ctx, qry, batchSize)
	if err != nil {
		return nil, fmt.Errorf("fetch and mark batch: %w", err)
	}
	defer rows.Close()

	var events []*model.Event
	for rows.Next() {
		var event model.Event
		if err := rows.Scan(&event.ID, &event.OrderID, &event.Status, &event.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}
		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return events, nil
}
