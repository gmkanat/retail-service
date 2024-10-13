package service

import (
	"context"
	"log"
	"sync"
	"time"

	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/infra"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/repository_raw/outbox"
)

type OutboxProcessor struct {
	repo          *outbox.Repository
	kafkaProducer *infra.KafkaProducer
	pollInterval  time.Duration
	stopCh        chan struct{}
	wg            sync.WaitGroup
}

func NewOutboxProcessor(repo *outbox.Repository, kafkaProducer *infra.KafkaProducer, pollInterval time.Duration) *OutboxProcessor {
	return &OutboxProcessor{
		repo:          repo,
		kafkaProducer: kafkaProducer,
		pollInterval:  pollInterval,
		stopCh:        make(chan struct{}),
	}
}

func (p *OutboxProcessor) Start(ctx context.Context) {
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
				p.processBatch(ctx)
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

func (p *OutboxProcessor) processBatch(ctx context.Context) {
	events, err := p.repo.FetchNextBatch(ctx, 10)
	if err != nil {
		log.Printf("error fetching events from outbox: %v", err)
		return
	}

	for _, event := range events {
		if err := p.kafkaProducer.SendEvent(ctx, event); err != nil {
			log.Printf("failed to send event to Kafka: %v", err)
			continue
		}

		if err := p.repo.MarkAsSent(ctx, []*model.Event{event}); err != nil {
			log.Printf("failed to mark event as sent: %v", err)
		}
	}
}
