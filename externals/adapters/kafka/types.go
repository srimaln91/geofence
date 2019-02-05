package kafka

import "github.com/Shopify/sarama"

type KClient struct {
	Producer             sarama.SyncProducer
	Consumer             sarama.Consumer
	Client               sarama.Client
	ConsumerGroupHandler ConsumerGroupHandler
	SchemaRegistry       *SchemaRegistry
}

type ConsumerGroupHandler struct {
	WriteChannel   chan Message
	SchemaRegistry *SchemaRegistry
}

type Message struct {
	SchemaId  int
	Topic     string
	Partition int32
	Offset    int64
	Key       string
	Value     string
}
