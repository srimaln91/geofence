package config

type schemaVersion struct {
	Topic   string
	Version int
}

type KafkaConfig struct {
	Brokers []string
	Schemas []schemaVersion
}
