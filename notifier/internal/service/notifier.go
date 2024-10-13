package service

import (
	"context"
	"github.com/IBM/sarama"
	"gitlab.ozon.dev/kanat_9999/homework/notifier/internal/config"
	"gitlab.ozon.dev/kanat_9999/homework/notifier/internal/infra"
	"log"
	"sync"
)

type Notifier struct {
	consumerGroup *infra.ConsumerGroup
}

func NewNotifier(cfg *config.Config) *Notifier {
	consumerGroup, err := infra.NewConsumerGroup(cfg)
	if err != nil {
		log.Fatalf("failed to create consumer group: %v", err)
	}

	return &Notifier{consumerGroup: consumerGroup}
}

func (n *Notifier) Start(ctx context.Context, cfg *config.Config, wg *sync.WaitGroup) {
	n.consumerGroup.Start(ctx, wg, cfg.KafkaTopic, &NotifierHandler{})
}

type NotifierHandler struct{}

func (h *NotifierHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *NotifierHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *NotifierHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("message received: topic=%s partition=%d offset=%d key=%s value=%s",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))

		session.MarkMessage(msg, "")
	}

	return nil
}
