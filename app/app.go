package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/srimaln91/geofence/config"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/srimaln91/geofence/api/controllers"
	"github.com/srimaln91/geofence/events"
	kclient "github.com/srimaln91/geofence/externals/adapters/kafka"
	"github.com/srimaln91/geofence/externals/adapters/mqtt"
	"github.com/srimaln91/geofence/externals/adapters/tile38"
	"github.com/srimaln91/geofence/externals/repositories"
	"github.com/srimaln91/geofence/mqtt_controller"
	"github.com/srimaln91/geofence/queue"
	ws "github.com/srimaln91/geofence/ws_controller"
)

// Init bootstraps the application
func Init(appConfig *config.AppConfig) {

	//Initializ tile38 connection
	t38Config := tile38.RedisConfig{
		DB:       appConfig.Tile38.DB,
		HostName: appConfig.Tile38.HostName,
		Password: appConfig.Tile38.Password,
		Port:     appConfig.Tile38.Port,
	}

	T38Adapter := tile38.Tile38Adapter{}
	T38Adapter.New(&t38Config)

	// Initialize Kafka
	schemaConfig := make(map[string]int)

	for _, schemas := range appConfig.Kafka.Schemas {
		schemaConfig[schemas.Topic] = schemas.Version
	}

	kafkaClient, err := kclient.New(appConfig.Kafka.Brokers, schemaConfig)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to kafka brokers: ", appConfig.Kafka.Brokers, appConfig.Kafka.Schemas)

	geofenceRepository := repositories.GeofenceRepository{
		T38Adapter: &T38Adapter,
	}

	// Initialize MQTT adapter and controller
	log.Println("Connecting to MQTT Broker: ", appConfig.Mqtt.Broker)
	mqttAdapter := mqtt.New(appConfig.Mqtt.Broker)
	mqttController := mqtt_controller.New(&geofenceRepository, &mqttAdapter)

	//Create Queue
	driverQueue := queue.ItemQueue{}
	driverQueue.New()

	wsController := ws.WSController{
		GeofenceRepository: &geofenceRepository,
	}

	mainController := controllers.MainController{
		GeofenceRepository: &geofenceRepository,
		WsController:       &wsController,
		DriverQueue:        &driverQueue,
	}

	eventController := events.EventController{
		GeofenceRepository: &geofenceRepository,
		WsController:       &wsController,
		KClient:            &kafkaClient,
		MqttController:     &mqttController,
		DriverQueue:        &driverQueue,
	}

	wsController.Init()

	// Listen for geofence trigers
	go eventController.SubscribeFence("myhook1")

	// Listen for location heartbeat data
	go mqttController.ListenHeartbeat()

	//Listen for geofence triggers
	go eventController.ListenFence()

	//Initialize router
	r := mux.NewRouter()

	// Routes consist of a path and a handler function.
	r.HandleFunc("/", mainController.BaseRoute)
	r.HandleFunc("/geofence", mainController.AddFence).Methods("POST")
	r.HandleFunc("/vehicle", mainController.AddVehicle).Methods("POST")
	r.HandleFunc("/queue", mainController.ReadQueue).Methods("GET")

	r.HandleFunc("/listen", mainController.Listen).Methods("POST")
	r.HandleFunc("/ws", wsController.HandleWSConn)

	//Enable CORS support
	c := cors.Default()

	handler := c.Handler(r)

	// Bind to a port and pass our router in
	log.Println("HTTP server is listening on: ", fmt.Sprintf("%s:%d", appConfig.Http.Host, appConfig.Http.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", appConfig.Http.Host, appConfig.Http.Port), handler))

}
