package kafka

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/gogovan-korea/ggx-kr-service-utils/logger"
	"github.com/gogovan-korea/ggx-kr-service-utils/saga/msg"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Worker func(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int)

type Consumer interface {
	msg.Consumer
	ConsumeTopic(ctx context.Context, groupTopics []string, numWorker int, worker Worker)
}

var DefaultAckWait = time.Second * 30

type consumerGroup struct {
	Brokers    []string
	GroupID    string
	logger     logger.Logger
	serializer Serializer
	ackWait    time.Duration
}

func NewConsumerGroup(brokers []string, groupID string, logger logger.Logger) Consumer {
	return &consumerGroup{Brokers: brokers, GroupID: groupID, logger: logger, serializer: DefaultSerializer, ackWait: DefaultAckWait}
}

//Listen in one topic only
func (c *consumerGroup) Listen(ctx context.Context, channel string, consumer msg.ReceiveMessageFunc) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: c.Brokers,
		GroupID: c.GroupID,
		Topic:   channel,
	})

	defer func(reader *kafka.Reader) {
		err := reader.Close()
		if err != nil {
			c.logger.Error("error closing kafka-go reader", zap.Error(err))
		}
	}(reader)

	for {
		err := c.receiveMessage(ctx, reader, consumer)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}
}

func (c *consumerGroup) receiveMessage(ctx context.Context, reader *kafka.Reader, consumer msg.ReceiveMessageFunc) error {
	m, err := reader.FetchMessage(ctx)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	var message msg.Message
	message, err = c.serializer.Deserialize(m)
	if err != nil {
		return err
	}

	wCtx, cancel := context.WithTimeout(ctx, c.ackWait)
	defer cancel()

	errc := make(chan error)
	go func() {
		errc <- consumer(wCtx, message)
	}()

	select {
	case err = <-errc:
		if err == nil {
			if ackErr := reader.CommitMessages(ctx, m); ackErr != nil {
				c.logger.Error("error acknowledging message", zap.Error(err))
			}
		}
	case <-ctx.Done():
		c.logger.Warn("listener has closed; in-progress message processing is terminated")
	case <-wCtx.Done():
		c.logger.Warn("timed out waiting for message consumer to finish")
	}

	return nil
}

func (c *consumerGroup) ConsumeTopic(ctx context.Context, groupTopics []string, numWorker int, worker Worker) {
	r := NewKafkaReader(c.Brokers, groupTopics, c.GroupID)

	defer func() {
		if err := r.Close(); err != nil {
			c.logger.Warn("Close consumer topic", zap.Error(err))
		}
	}()

	c.logger.Info("Starting consumer topic",
		zap.String("GroupID", c.GroupID),
		zap.Any("groupTopics", groupTopics),
		zap.Int("numWorker", numWorker))

	wg := &sync.WaitGroup{}
	for i := 0; i < numWorker; i++ {
		wg.Add(1)
		go worker(ctx, r, wg, i)
	}
	wg.Wait()
}
func (c *consumerGroup) Close(context.Context) error {
	c.logger.Warn("closing message source")
	return nil
}
