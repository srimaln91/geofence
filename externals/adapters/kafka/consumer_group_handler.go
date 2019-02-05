package kafka

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/linkedin/goavro"

	"github.com/Shopify/sarama"
)

// Implementing ConsumerGroup interface

//Setup sets up the group
func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

// Cleanup gracefully terminates the group
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim listen for messages
func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)

		//Decode message
		schemaId := binary.BigEndian.Uint32(msg.Value[1:5])

		log.Println("Decode Schema ID: ", schemaId)

		codec, err := goavro.NewCodec(h.SchemaRegistry.Schemas[msg.Topic][int(schemaId)].Schema)
		if err != nil {
			log.Println("Codec: ", err)
		}

		// Convert binary Avro data back to native Go form
		native, _, err := codec.NativeFromBinary(msg.Value[5:])
		if err != nil {
			return err
		}

		// Convert native Go form to textual Avro data
		textual, err := codec.TextualFromNative(nil, native)

		if err != nil {
			return err
		}
		message := Message{int(schemaId), msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(textual)}

		h.WriteChannel <- message
		sess.MarkMessage(msg, "")
	}
	return nil
}
