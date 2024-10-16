package service

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/transaction"
	"log"
	"sync"
	"time"

	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/infra"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/repository_raw/outbox"
)

type OutboxProcessor struct {
	repo          *outbox.Repository
	kafkaProducer *infra.KafkaProducer
	pollInterval  time.Duration
	stopCh        chan struct{}
	wg            sync.WaitGroup
	txManager     *transaction.TxManager
}

func NewOutboxProcessor(
	repo *outbox.Repository,
	kafkaProducer *infra.KafkaProducer,
	pollInterval time.Duration,
	txManager *transaction.TxManager,
) *OutboxProcessor {
	return &OutboxProcessor{
		repo:          repo,
		kafkaProducer: kafkaProducer,
		pollInterval:  pollInterval,
		stopCh:        make(chan struct{}),
		txManager:     txManager,
	}
}

func (p *OutboxProcessor) Start(ctx context.Context, batchSize int) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("context cancelled, stopping outbox processor...")
				return
			case <-p.stopCh:
				log.Println("stop signal received, stopping outbox processor...")
				return
			default:
				if err := p.processBatch(ctx, batchSize); err != nil {
					log.Printf("failed to process batch: %v", err)
				}
				time.Sleep(p.pollInterval)
			}
		}
	}()
}

func (p *OutboxProcessor) Stop() {
	close(p.stopCh)
	p.wg.Wait()
	log.Println("outbox processor stopped")
}

func (p *OutboxProcessor) processBatch(ctx context.Context, batchSize int) error {
	return p.txManager.WithRepeatableReadTx(ctx, func(c context.Context) error {
		events, err := p.repo.FetchAndMarkBatch(c, batchSize)
		if err != nil {
			return fmt.Errorf("fetch and mark batch: %w", err)
		}

		if len(events) == 0 {
			return nil
		}

		for _, event := range events {
			if err := p.kafkaProducer.SendEvent(c, event); err != nil {
				return fmt.Errorf("send event to Kafka: %w", err)
			}
		}
		return nil
	})
}
