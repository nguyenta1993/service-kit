package kafka

type Config struct {
	Config *ConfigDetail
	Topics []TopicConfig
}

type ConfigDetail struct {
	Brokers    []string
	GroupID    string
	InitTopics bool
	NumWorker  int
}

type TopicConfig struct {
	TopicName         string
	NumPartitions     int
	ReplicationFactor int
}

type ConsumerConfig struct {
	Topics   []string
	PoolSize int
	Worker   Worker
}
