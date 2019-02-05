package events

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	kclient "github.com/srimaln91/geofence/externals/adapters/kafka"
	"github.com/srimaln91/geofence/externals/adapters/tile38"
	"github.com/srimaln91/geofence/externals/repositories"
	"github.com/srimaln91/geofence/mqtt_controller"
	"github.com/srimaln91/geofence/queue"
	ws "github.com/srimaln91/geofence/ws_controller"
)

type EventController struct {
	WsController       *ws.WSController
	KClient            *kclient.KClient
	MqttController     *mqtt_controller.MqttController
	GeofenceRepository *repositories.GeofenceRepository
	DriverQueue        *queue.ItemQueue
}

type FenceHook struct {
	Command string    `json:"command"`
	Group   string    `json:"group"`
	Detect  string    `json:"detect"`
	Hook    string    `json:"hook"`
	Key     string    `json:"key"`
	Time    time.Time `json:"time"`
	ID      string    `json:"id"`
	Object  struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"object"`
}

var fenceData = make(chan *tile38.SubMessage)

const (
	topic = "geofence_trigger"
)

// New creates a new kafka controller
func New(brokers []string, wsController *ws.WSController, mqttController *mqtt_controller.MqttController, kClient *kclient.KClient) *EventController {

	return &EventController{
		WsController:   wsController,
		MqttController: mqttController,
		KClient:        kClient,
	}

}

// ListenFence listens for geofence channel
func (event *EventController) ListenFence() {

	go event.ConsumeKafka()

	for r := range fenceData {

		part, offset, err := event.KClient.ProduceMessage(topic, []byte(r.Payload))

		if err != nil {
			fmt.Println("Produce: ", err)
		}

		fmt.Println("Produced | Offset:", offset, " | Partition: ", part)

	}

}

// ConsumerKafka consumes messages from kafka topic
func (event *EventController) ConsumeKafka() {

	msgChan := event.KClient.GetMessageChannel()

	go func() {
		for msg := range msgChan {

			encMessage := string(msg.Value)
			fmt.Println(encMessage)

			fenceDecoded := FenceHook{}

			err := json.Unmarshal([]byte(msg.Value), &fenceDecoded)

			if err != nil {
				log.Println(err)
			}

			vehicle := queue.Item{
				Id: fenceDecoded.ID,
			}

			if fenceDecoded.Detect == "enter" {

				// Add/Remove items from queue
				event.DriverQueue.Enqueue(vehicle)

			} else {

				// if event.DriverQueue.IsEmpty() == false {
				event.DriverQueue.Remove(vehicle)
				// }

			}

			fmt.Println(event.DriverQueue.GetItems())

			// event.WsController.BroadcastMessage(msg.Value)
			event.MqttController.FireMqttFenceNotification("fence/"+fenceDecoded.ID, encMessage)

		}
	}()

	event.KClient.ConsumeMessages(topic, msgChan)

}

func (event *EventController) SubscribeFence(fenceName string) {

	// Subscribe to channel data
	go event.GeofenceRepository.Subscribe(fenceData, fenceName)

}
