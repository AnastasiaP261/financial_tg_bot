package consumer

import (
	"context"
	"log"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/kafka"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/report"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/utils/logs"
	"go.uber.org/zap"
)

type defaultReportsHandler interface {
	CreateReport(ctx context.Context, rawReq string) (report.CreateReportResponse, error)
}

// Consumer represents a Sarama consumer group consumer.
type Consumer struct {
	defaultReportsHandler defaultReportsHandler
}

func StartConsumerGroup(ctx context.Context, defaultReportsHandler defaultReportsHandler) error {
	consumerGroupHandler := Consumer{
		defaultReportsHandler: defaultReportsHandler,
	}

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
	for msg := range claim.Messages() {
		logs.Info(
			"New message received",
			zap.String("from topic", msg.Topic),
			zap.Int64("offset", msg.Offset),
			zap.Int32("partition", msg.Partition),
			zap.String("key", string(msg.Key)),
			zap.String("value", string(msg.Value)),
		)

		switch string(msg.Key) {
		case "get_report":
			_, err := c.defaultReportsHandler.CreateReport(context.Background(), string(msg.Value))
			if err != nil {
				logs.Error(
					"handle msg error",
					zap.Error(err),
					zap.String("key", string(msg.Key)),
					zap.String("value", string(msg.Value)),
				)
				return errors.Wrap(err, "defaultReportsHandler.CreateReport")
			}
		default:
			logs.Error("read invalid message")
		}
		session.MarkMessage(msg, "")
	}

	return nil
}
