package consumer

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/kafka"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/logs"
	"go.uber.org/zap"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

func StartConsumerGroup(ctx context.Context) error {
	consumerGroupHandler := Consumer{}

	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	switch kafka.Assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategySticky}
	case "round-robin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRange}
	default:
		log.Fatalf("Unrecognized consumer group partition assignor: %s", kafka.Assignor)
	}

	// Create consumer group
	consumerGroup, err := sarama.NewConsumerGroup(kafka.BrokersList, kafka.ConsumerGroup, config)
	if err != nil {
		return errors.Wrap(err, "starting consumer group")
	}

	err = consumerGroup.Consume(ctx, []string{kafka.Topic}, &consumerGroupHandler)
	if err != nil {
		return errors.Wrap(err, "consuming via handler")
	}
	return nil
}

// Consumer represents a Sarama consumer group consumer.
type Consumer struct{}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	logs.Info("c - setup")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	logs.Info("c - cleanup")
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		printMessage(message)
		session.MarkMessage(message, "")
	}

	return nil
}

func printMessage(msg *sarama.ConsumerMessage) {
	logs.Info(
		"New message received",
		zap.String("from topic", msg.Topic),
		zap.Int64("offset", msg.Offset),
		zap.Int32("partition", msg.Partition),
		zap.String("key", string(msg.Key)),
		zap.String("value", string(msg.Value)),
	)

	// Emulate Work loads
	time.Sleep(1 * time.Second)

	logs.Info(
		"Successful to read message",
		zap.String("value", string(msg.Value)),
	)
}
