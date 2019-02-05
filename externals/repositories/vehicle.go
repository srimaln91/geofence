package repositories

import (
	"github.com/srimaln91/geofence/externals/adapters/tile38"
)

type Vehicle struct {
	Id         string
	Collection string
	Location   struct {
		Long string
		Lat  string
	}
}

type VehicleRepository struct {
	T38Adapter *tile38.Tile38Adapter
}

// AddFence adds a new geofence
func (repository GeofenceRepository) SetLocation(vehicle Vehicle) (interface{}, error) {

	command := tile38.Command{
		IsWrite:    true,
		QueryParts: make([]interface{}, 0),
	}

	// New Geofence
	command.AddParam("SET")
	command.AddParam(vehicle.Collection)
	command.AddParam(string(vehicle.Id))

	command.AddParam("POINT")
	command.AddParam(vehicle.Location.Lat)
	command.AddParam(vehicle.Location.Long)

	return repository.T38Adapter.RunCommand(&command)
}
