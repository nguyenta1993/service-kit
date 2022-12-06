package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"strconv"

	"github.com/gogovan-korea/ggx-kr-service-utils/logger"
	"github.com/segmentio/kafka-go/sasl/plain"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var dialer *kafka.Dialer

func UseKafka(ctx context.Context, logger logger.Logger, cfg *Config, consumerConfig *ConsumerConfig) *kafka.Conn {
	kafkaConn := connectKafkaBrokers(ctx, logger, cfg)

	if cfg.Config.InitTopics {
		initKafkaTopics(ctx, logger, cfg, kafkaConn)
	}

	if consumerConfig != nil {
		consumerTopics(ctx, logger, cfg, consumerConfig)
	}

	return kafkaConn
}

func connectKafkaBrokers(ctx context.Context, logger logger.Logger, cfg *Config) *kafka.Conn {
	InitDialer(&cfg.Dialer)

	kafkaConn, err := dialer.DialContext(ctx, "tcp", cfg.Config.Brokers[0])
	if err != nil {
		logger.Error("kafka connection", zap.Error(err))
	}

	brokers, err := kafkaConn.Brokers()
	if err != nil {
		logger.Error("kafkaConn.Brokers", zap.Error(err))
		return nil
	}
	logger.Info("kafka connected to brokers", zap.Any("Brokers", brokers))
	return kafkaConn
}

func InitDialer(config ...*DialerConfig) {

	if dialer != nil {
		return
	}
	if len(config) == 0 {
		dialer = kafka.DefaultDialer
	}
	cfg := config[0]
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	dialer = &kafka.Dialer{
		DualStack: true,
		SASLMechanism: plain.Mechanism{
			Username: cfg.Username, // access key
			Password: cfg.Password, // secret
		},
		TLS: &tls.Config{
			InsecureSkipVerify: true,
			RootCAs:            rootCAs,
		},
	}
}

func initKafkaTopics(ctx context.Context, logger logger.Logger, cfg *Config, kafkaConn *kafka.Conn) {
	controller, err := kafkaConn.Controller()
	if err != nil {
		logger.Error("kafkaConn.Controller", zap.Error(err))
		return
	}

	controllerURI := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	logger.Info("kafka controller uri", zap.String("controllerURI", controllerURI))

	conn, err := dialer.DialContext(ctx, "tcp", controllerURI)
	if err != nil {
		logger.Error("DialContext controller", zap.Error(err))
		return
	}
	defer conn.Close() // nolint: errcheck

	var topics []kafka.TopicConfig

	for _, topic := range cfg.Topics {
		topics = append(topics, kafka.TopicConfig{
			Topic:             topic.TopicName,
			NumPartitions:     topic.NumPartitions,
			ReplicationFactor: topic.ReplicationFactor,
		})
	}

	if err := conn.CreateTopics(topics...); err != nil {
		logger.Error("CreateTopics", zap.Error(err))
		return
	}
}

func consumerTopics(ctx context.Context, logger logger.Logger, cfg *Config, consumerConfig *ConsumerConfig) {
	cg := NewConsumerGroup(cfg.Config.Brokers, cfg.Config.GroupID, logger)
	go cg.ConsumeTopic(ctx, consumerConfig.Topics, consumerConfig.PoolSize, consumerConfig.Worker)
}
