package kafka

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/linkedin/goavro"
)

// ProduceMessage produces a message
func (kc *KClient) ProduceMessage(topic string, message []byte) (int32, int64, error) {

	var partition int32 = -1

	//Get default version
	defaultVersion := kc.SchemaRegistry.DefaultVersions[topic]
	schemaVersion := kc.SchemaRegistry.Schemas[topic][defaultVersion].Version

	binarySchemaId := make([]byte, 4)
	binary.BigEndian.PutUint32(binarySchemaId, uint32(schemaVersion))

	log.Println("Binary Schema ID: ", binarySchemaId)

	//Encode the message
	schema := kc.SchemaRegistry.Schemas[topic][schemaVersion].Schema
	log.Println("Producer encoding Schema: ", schema)
	codec, err := goavro.NewCodec(schema)
	if err != nil {
		fmt.Println("Schema: ", err)
	}

	encodedMessage, _, err := codec.NativeFromTextual(message)
	if err != nil {
		fmt.Println("Encode", err)
	}

	// Convert native Go form to binary Avro data
	binaryValue, err := codec.BinaryFromNative(nil, encodedMessage)
	if err != nil {
		return 0, 0, err
	}

	var binaryMsg []byte
	// first byte is magic byte, always 0 for now
	binaryMsg = append(binaryMsg, byte(0))
	//4-byte schema ID as returned by the Schema Registry
	binaryMsg = append(binaryMsg, binarySchemaId...)
	//avro serialized data in Avro’s binary encoding
	binaryMsg = append(binaryMsg, binaryValue...)

	// fmt.Println(encodedMessage)
	praparedMessage := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: partition,
		Value:     sarama.StringEncoder(binaryMsg),
	}

	return kc.Producer.SendMessage(praparedMessage)

}
