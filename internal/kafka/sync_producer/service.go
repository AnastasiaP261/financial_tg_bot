package sync_producer

import (
	"fmt"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/kafka"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/logs"
	"go.uber.org/zap"
	"time"

	"github.com/Shopify/sarama"
)

type Service struct {
	producer sarama.SyncProducer
}

func New() (*Service, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	// Waits for all in-sync replicas to commit before responding.
	config.Producer.RequiredAcks = sarama.WaitForAll
	// The total number of times to retry sending a message (default 3).
	config.Producer.Retry.Max = 3
	// How long to wait for the cluster to settle between retries (default 100ms).
	config.Producer.Retry.Backoff = time.Millisecond * 250
	// idempotent producer has a unique producer ID and uses sequence IDs for each message,
	// allowing the broker to ensure, on a per-partition basis, that it is committing ordered messages with no duplication.
	//config.Producer.Idempotent = true
	if config.Producer.Idempotent == true {
		config.Producer.Retry.Max = 1
		config.Net.MaxOpenRequests = 1
	}
	//  Successfully delivered messages will be returned on the Successes channel
	config.Producer.Return.Successes = true
	// Generates partitioners for choosing the partition to send messages to (defaults to hashing the message key)
	_ = config.Producer.Partitioner

	producer, err := sarama.NewSyncProducer(kafka.BrokersList, config)
	if err != nil {
		return nil, fmt.Errorf("starting Sarama producer: %w", err)
	}

	return &Service{producer: producer}, nil
}

func (s *Service) Close() error {
	return s.producer.Close()
}

func (s *Service) SendNewMsg(key, value string) {
	// Inject info into message
	msg := sarama.ProducerMessage{
		Topic: kafka.Topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	p, o, err := s.producer.SendMessage(&msg)
	if err != nil {
		logs.Error("Failed to send message", zap.Error(err))
	}
	logs.Info(
		"Successful to write message",
		zap.String("topic", kafka.Topic),
		zap.Int64("offset", o),
		zap.Int32("partition", p),
	)
}
