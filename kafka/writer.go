package kafka

import (
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

const MaxAttempts = 3

func NewWriter(writer *kafka.Writer) *kafka.Writer {
	if writer.Balancer == nil {
		writer.Balancer = &kafka.LeastBytes{}
	}

	if writer.RequiredAcks == kafka.RequireNone {
		writer.RequiredAcks = kafka.RequireAll
	}

	writer.MaxAttempts = MaxAttempts
	writer.Compression = compress.Snappy

	return writer
}
