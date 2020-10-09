package db

import (
	"Edgex-Export_Receiver/app/db/Mongo"
	"Edgex-Export_Receiver/app/domain"
)

type DeviceRepos interface {
	SelectAll()	([]domain.Device, error)
	Select(id string) (domain.Device, error)
	SelectByName(name string) (domain.Device, error)
	SelectByEdgexId(edgexId string) (domain.Device, error)
	Insert(device *domain.Device) (string, error)
	Delete(id string) error
	DeleteByEdgexId(edgexId string) error
	Update(device *domain.Device) (domain.Device, error)
}

func GetDeviceRepos() DeviceRepos {
	if Mongo.DS.S != nil {
		return DeviceRepos(&Mongo.DeviceMongoRepository{})
	} else{
		return nil
	}
}