package infra

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"gitlab.ozon.dev/kanat_9999/homework/loms/internal/model"
	"log"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaProducer(brokers []string, topic string) (*KafkaProducer, error) {
	config := prepareConfig()
	log.Printf("Kafka brokers: %v", brokers)
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (kp *KafkaProducer) SendEvent(ctx context.Context, event *model.Event) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: kp.topic,
		Value: sarama.ByteEncoder(eventBytes),
	}

	_, _, err = kp.producer.SendMessage(msg)
	if err != nil {
		log.Printf("failed to send message to Kafka: %v", err)
		return err
	}

	log.Printf("sent event to Kafka: %v", event)
	return nil
}

func (kp *KafkaProducer) Close() error {
	return kp.producer.Close()
}

func prepareConfig() *sarama.Config {
	c := sarama.NewConfig()
	c.Producer.Partitioner = sarama.NewHashPartitioner // partition by key
	c.Producer.RequiredAcks = sarama.WaitForAll        // guarentee delivery
	c.Net.MaxOpenRequests = 1                          // guarentee strict order
	c.Producer.Return.Successes = true
	c.Producer.Return.Errors = true
	return c
}
