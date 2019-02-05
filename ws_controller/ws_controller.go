package ws_controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/xid"
	"github.com/srimaln91/geofence/externals/repositories"
)

type WSController struct {
	Clients            map[string]*websocket.Conn
	GeofenceRepository *repositories.GeofenceRepository
}

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// VehicleLocationRequest struct for the vehicle location
type VehicleLocationRequest struct {
	ID         string   `json:"id"`
	Collection string   `json:"collection"`
	Location   Location `json:"location"`
}

// Location location struct
type Location struct {
	Longitude string `json:"long"`
	Latitude  string `json:"lat"`
}

var locationheartbeat = make(chan VehicleLocationRequest)

// Init initializes the controller
func (ws *WSController) Init() {
	ws.Clients = make(map[string]*websocket.Conn)
}

// HandleWSConn handles the websocket connection
func (ws *WSController) HandleWSConn(w http.ResponseWriter, r *http.Request) {

	socket, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
	}

	defer socket.Close()

	guid := xid.New().String()

	ws.Clients[guid] = socket

	for {
		var msg VehicleLocationRequest

		err := socket.ReadJSON(&msg)

		if err != nil {
			delete(ws.Clients, guid)
			break
		}

		locationheartbeat <- msg
	}
}

// ConsumeHeartbeat consumes heartbeat data
func (ws *WSController) ConsumeHeartbeat() {

	for message := range locationheartbeat {

		// Create and populate Vehicle struct
		vehicle := repositories.Vehicle{}

		vehicle.Collection = message.Collection
		vehicle.Id = message.ID
		vehicle.Location.Lat = message.Location.Latitude
		vehicle.Location.Long = message.Location.Longitude

		_, err := ws.GeofenceRepository.SetLocation(vehicle)

		if err != nil {
			fmt.Println(err.Error())
		}

	}
}

func (ws *WSController) BroadcastMessage(message []byte) {

	for _, wsClient := range ws.Clients {
		wsClient.WriteMessage(1, message)
	}

}
