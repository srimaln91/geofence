package config

type AppConfig struct {
	Kafka  *KafkaConfig
	Http   *HttpConfig
	Mqtt   *MqttConfig
	Tile38 *Tile38Config
}
