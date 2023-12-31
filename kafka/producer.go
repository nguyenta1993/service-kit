package kafka

import (
	"context"
	"time"

	"github.com/nguyenta1993/service-kit/logger"
	"github.com/nguyenta1993/service-kit/saga/msg"
	"go.uber.org/zap"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	log        logger.Logger
	w          *kafka.Writer
	serializer Serializer
}

func NewProducer(log logger.Logger, brokers []string, async bool, timeout time.Duration, cfg ...*DialerConfig) *Producer {
	InitDialer(cfg...)
	w := kafka.NewWriter(kafka.WriterConfig{
		Dialer:       dialer,
		Async:        async,
		BatchTimeout: timeout,
		Brokers:      brokers,
	})
	return &Producer{log: log, w: NewWriter(w), serializer: DefaultSerializer}
}

func (p *Producer) PublishMessage(ctx context.Context, msgs ...kafka.Message) error {
	return p.w.WriteMessages(ctx, msgs...)
}

func (p *Producer) Close(context.Context) error {
	p.log.Info("closing message destination")
	err := p.w.Close()
	if err != nil {
		p.log.Error("error closing message destination", zap.Error(err))
	}
	return err
}

// Send to one topic only
func (p *Producer) Send(ctx context.Context, channel string, message msg.Message) error {
	kafkaMsg, err := p.serializer.Serialize(message)
	if err != nil {
		p.log.Error("failed to marshal message", zap.Error(err))
		return err
	}

	kafkaMsg.Topic = channel

	return p.w.WriteMessages(ctx, kafkaMsg)
}
