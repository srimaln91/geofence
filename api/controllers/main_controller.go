package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/srimaln91/geofence/queue"

	"github.com/gorilla/websocket"
	"github.com/srimaln91/geofence/api/handlers"
	"github.com/srimaln91/geofence/externals/adapters/tile38"
	"github.com/srimaln91/geofence/externals/repositories"
	ws "github.com/srimaln91/geofence/ws_controller"
)

//FenceRequest data struct
type FenceRequest struct {
	ID         string          `json:"id"`
	Collection string          `json:"collection"`
	GeoJSON    json.RawMessage `json:"geojson"`
}

// Location location struct
type Location struct {
	Longitude string `json:"long"`
	Latitude  string `json:"lat"`
}

// VehicleLocationRequest struct for the vehicle location
type VehicleLocationRequest struct {
	ID         string   `json:"id"`
	Collection string   `json:"collection"`
	Location   Location `json:"location"`
}

// MainController struct for the controler
type MainController struct {
	GeofenceRepository *repositories.GeofenceRepository
	WsController       *ws.WSController
	DriverQueue        *queue.ItemQueue
}

var clients = make(map[*websocket.Conn]bool) // connected clients
// var broadcast = make(chan Message)           // broadcast channel

var locationheartbeat = make(chan VehicleLocationRequest)

var fenceData = make(chan *tile38.SubMessage)

// BaseRoute is the initial route
func (c *MainController) BaseRoute(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Hello There. Welcome to API")
}

// AddFence adds a new geofence
func (c *MainController) AddFence(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	fenceRequest := FenceRequest{}

	err := json.NewDecoder(r.Body).Decode(&fenceRequest)

	response := handlers.Response{}
	response.Error = true

	if err != nil {
		response.Message = err.Error()
		response.ResponseJSON(w)
		return
	}

	fence := repositories.GeoFence{
		ID:         fenceRequest.ID,
		Name:       "myhook1",
		Collection: fenceRequest.Collection,
		GeoJSON:    string(fenceRequest.GeoJSON),
	}

	result, err := c.GeofenceRepository.AddFence(fence)

	if err != nil {
		response.Message = err.Error()
		response.ReponseError(w)
		return
	}

	// Subscribe to channel data
	// go c.GeofenceRepository.Subscribe(fenceData, fence.Name)

	response.Error = false
	response.Data = result
	response.ResponseJSON(w)

}

// AddVehicle adds a new vehicle
func (c *MainController) AddVehicle(w http.ResponseWriter, r *http.Request) {

	locationRequest := VehicleLocationRequest{}

	err := json.NewDecoder(r.Body).Decode(&locationRequest)

	response := handlers.Response{}
	response.Error = true

	if err != nil {
		response.Message = err.Error()
		return
	}

	// Create and populate Vehicle struct
	vehicle := repositories.Vehicle{}

	vehicle.Collection = locationRequest.Collection
	vehicle.Id = locationRequest.ID
	vehicle.Location.Lat = locationRequest.Location.Latitude
	vehicle.Location.Long = locationRequest.Location.Longitude

	result, err := c.GeofenceRepository.SetLocation(vehicle)

	if err != nil {
		response.Message = err.Error()
		response.ReponseError(w)
		return
	}

	response.Error = false
	response.Data = result
	response.ResponseJSON(w)

}

func (c *MainController) ReadQueue(w http.ResponseWriter, r *http.Request) {

	response := handlers.Response{}

	type vehicle struct {
		ID string `json:"ID"`
	}
	type queueresp struct {
		Vehicle []queue.Item `json:"vehicle"`
	}

	preparedResp := queueresp{
		Vehicle: c.DriverQueue.GetItems(),
	}

	// if err != nil {
	// 	response.Message = err.Error()
	// 	response.ReponseError(w)
	// 	return
	// }

	response.Error = false
	response.Data = preparedResp
	response.ResponseJSON(w)

}

// Listen listens for Tile38 webhooks
func (c *MainController) Listen(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	reqData, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(reqData))

	//Send geofence data to the client
	for wsClient := range clients {
		wsClient.WriteMessage(1, reqData)
	}

	w.WriteHeader(200)
	return

}
