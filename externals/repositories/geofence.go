package repositories

import (
	"github.com/srimaln91/geofence/externals/adapters/tile38"
)

type GeoFence struct {
	ID         string
	Name       string
	Collection string
	GeoJSON    string
}

type GeofenceRepository struct {
	T38Adapter *tile38.Tile38Adapter
}

// AddFence adds a new geofence
func (repository GeofenceRepository) AddFence(fence GeoFence) (interface{}, error) {

	command := tile38.Command{
		IsWrite:    true,
		QueryParts: make([]interface{}, 0),
	}

	// New Geofence
	command.AddParam("SETCHAN")
	command.AddParam(fence.Name)
	command.AddParam("INTERSECTS")
	command.AddParam(fence.Collection)
	command.AddParam("FENCE")
	command.AddParam("DETECT")
	command.AddParam("enter,exit")
	command.AddParam("OBJECT")
	command.AddParam(fence.GeoJSON)

	return repository.T38Adapter.RunCommand(&command)
}

func (repository GeofenceRepository) Subscribe(receiver chan *tile38.SubMessage, channel string) {

	repository.T38Adapter.SubscribeChannel(receiver, channel)

}
