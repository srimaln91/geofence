package kafka

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
)

// ConsumePartitionMessages returns a channel that can be used to listen for messages
func (kc *KClient) ConsumePartitionMessages(topic string, msgChan chan *sarama.ConsumerMessage) {

	partitionList, _ := kc.Consumer.Partitions(topic)

	initialOffset := sarama.OffsetOldest

	for _, partition := range partitionList {

		pc, _ := kc.Consumer.ConsumePartition(topic, partition, initialOffset)

		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				msgChan <- message
			}
		}(pc)

	}
}

// ConsumeMessages consume messages as a consumer group
func (kc *KClient) ConsumeMessages(topic string, msgChan chan Message) {

	defer func() { _ = kc.Client.Close() }()

	// Start a new consumer group
	group, err := sarama.NewConsumerGroupFromClient("my-group", kc.Client)

	if err != nil {
		panic(err)
	}

	defer func() { _ = group.Close() }()

	// Track errors
	go func() {
		for err := range group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	// Iterate over consumer sessions.
	ctx := context.Background()

	for {
		topics := []string{topic}

		handler := ConsumerGroupHandler{
			WriteChannel:   msgChan,
			SchemaRegistry: kc.SchemaRegistry,
		}

		err := group.Consume(ctx, topics, handler)

		if err != nil {
			panic(err)
		}
	}

}

// GetMessageChannel returns a channel that can be used to receive consumes messages
func (kc *KClient) GetMessageChannel() chan Message {
	return make(chan Message, 256)
}
