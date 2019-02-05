package kafka

import (
	"github.com/Shopify/sarama"
)

// New populates KClient struct
func New(brokers []string, schemas map[string]int) (KClient, error) {

	//setup relevant config info
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Consumer.Return.Errors = true
	config.Version = sarama.V0_10_2_0

	// Register a producer
	producer, err := sarama.NewSyncProducer(brokers, config)

	if err != nil {
		return KClient{}, err
	}

	//Register a consumer
	consumer, err := sarama.NewConsumer(brokers, nil)

	if err != nil {
		return KClient{}, err
	}

	// Register Client
	client, err := sarama.NewClient(brokers, config)

	// Initialize schemas
	SchemaClient, err := GetSchemas(schemas)

	if err != nil {
		return KClient{}, err
	}

	return KClient{
		Producer:       producer,
		Consumer:       consumer,
		Client:         client,
		SchemaRegistry: SchemaClient,
	}, nil
}
