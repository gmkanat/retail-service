package infra

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/notifier/internal/config"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type ConsumerGroup struct {
	group sarama.ConsumerGroup
}

func NewConsumerGroup(cfg *config.Config) (*ConsumerGroup, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Session.Timeout = 10 * time.Second
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup(cfg.KafkaBrokers, "notifier-group", saramaConfig)
	if err != nil {
		return nil, err
	}

	return &ConsumerGroup{group: group}, nil
}

func (cg *ConsumerGroup) Start(ctx context.Context, wg *sync.WaitGroup, topic string, handler sarama.ConsumerGroupHandler) {
	defer wg.Done()
	for {
		if err := cg.group.Consume(ctx, []string{topic}, handler); err != nil {
			log.Printf("error from consumer: %v", err)
			return
		}

		if ctx.Err() != nil {
			return
		}
	}
}
