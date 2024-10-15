package outbox

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/infra"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/pgcluster"
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

func (n *Repository) CreateEvent(ctx context.Context, tx pgx.Tx, event model.Event) error {
	qry := `INSERT INTO outbox.notifier (order_id, status)
			VALUES ($1, $2);`

	if _, err := tx.Exec(ctx, qry, event.OrderID, event.Status); err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

func (n *Repository) FetchNextBatch(ctx context.Context, batchSize int) ([]*model.Event, error) {
	reader, err := n.cluster.GetWriter(ctx)
	if err != nil {
		return nil, err
	}

	qry := `SELECT id, order_id, status, created_at
			 FROM outbox.notifier
			 WHERE is_sent = false
			 ORDER BY id
			 LIMIT $1 FOR UPDATE SKIP LOCKED`

	rows, err := reader.Query(ctx, qry, batchSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*model.Event
	for rows.Next() {
		var event model.Event
		if err := rows.Scan(&event.ID, &event.OrderID, &event.Status, &event.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (n *Repository) MarkAsSent(ctx context.Context, events []*model.Event) error {
	if len(events) == 0 {
		return nil
	}

	ids := make([]int64, len(events))
	for i, event := range events {
		ids[i] = event.ID
	}

	writer, err := n.cluster.GetWriter(ctx)
	if err != nil {
		return err
	}

	qry := `UPDATE outbox.notifier
	        SET is_sent = true
	        WHERE id = ANY($1)`

	_, err = writer.Exec(ctx, qry, ids)
	return err
}
