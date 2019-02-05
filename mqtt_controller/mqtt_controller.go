package mqtt_controller

import (
	"encoding/json"
	"fmt"

	"github.com/srimaln91/geofence/externals/adapters/mqtt"
	"github.com/srimaln91/geofence/externals/repositories"
)

type MqttController struct {
	GeofenceRepository *repositories.GeofenceRepository
	MqttAdapter        *mqtt.MQTTClient
}

func New(GeofenceRepository *repositories.GeofenceRepository, mqttAdapter *mqtt.MQTTClient) MqttController {
	return MqttController{
		GeofenceRepository: GeofenceRepository,
		MqttAdapter:        mqttAdapter,
	}
}

func (mqc *MqttController) ListenHeartbeat() {

	msgChan := mqc.MqttAdapter.GetMessageChannel()

	mqc.MqttAdapter.SubscribeMessage("location/#", msgChan)

	go func() {
		for msg := range msgChan {
			// fmt.Println(msg.Payload())

			message := msg.Payload()
			vehicle := repositories.Vehicle{}

			uMarshallError := json.Unmarshal(message, &vehicle)

			if uMarshallError != nil {
				fmt.Println(uMarshallError)
			}
			_, err := mqc.GeofenceRepository.SetLocation(vehicle)

			if err != nil {
				fmt.Println(err.Error())
			}

		}
	}()

}

func (mqc *MqttController) FireMqttFenceNotification(topic string, message string) {
	mqc.MqttAdapter.PublishMessage(topic, message)
}
