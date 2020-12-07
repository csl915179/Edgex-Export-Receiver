package Mongo

import (
	"Edgex-Export_Receiver/app/domain"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type EventMongoRepository struct {
}

func (ar *EventMongoRepository) SelectNumber(number int64) ([]domain.Event, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(eventScheme)
	result := make([]domain.Event, 0)
	err := coll.Find(nil).Sort().All(&result)
	if err != nil {
		log.Println("Find Event failed !" + err.Error())
		return result, err
	}
	return result, err
}

func (ar *EventMongoRepository) Select(id string) (domain.Event, error){
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(eventScheme)
	result := domain.Event{}
	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
	}
	for  i:=0; i<len(result.Devices); i++ {
		coll := ds.S.DB(database).C(deviceScheme)
		PhysicalDevice := domain.Device{}
		coll.Find(bson.M{"_id": bson.ObjectIdHex(result.Devices[i].Id)}).One(&PhysicalDevice)
		result.Devices[i].AvailCpu = PhysicalDevice.Cpu - PhysicalDevice.CPUUsed
		result.Devices[i].AvailMem = PhysicalDevice.Memory - PhysicalDevice.MemoryUsed
		result.Devices[i].AvailDisk = PhysicalDevice.Disk - PhysicalDevice.DiskUsed
		result.Devices[i].AvailNetRate = PhysicalDevice.NetRate - PhysicalDevice.NetRateUsed
	}
	return result, nil
}

func (ar *EventMongoRepository) Insert(event *domain.Event) (string, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(eventScheme)
	err := coll.Insert(event)
	if err != nil {
		log.Println("Insert event failed !")
		return "", err
	}
	return event.Id.Hex(), nil
}

func (ar *EventMongoRepository) ExtractAll() ([]domain.Event, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(eventScheme)
	result := make([]domain.Event, 0)
	err := coll.Find(nil).All(&result)
	if err != nil {
		log.Println("Find all event failed !")
		return result, err
	}
	for i:=0; i<len(result); i++ {
		for  j:=0; j<len(result[i].Devices); j++ {
			coll := ds.S.DB(database).C(deviceScheme)
			PhysicalDevice := domain.Device{}
			coll.Find(bson.M{"_id": bson.ObjectIdHex(result[i].Devices[j].Id)}).One(&PhysicalDevice)
			result[i].Devices[j].AvailCpu = PhysicalDevice.Cpu - PhysicalDevice.CPUUsed
			result[i].Devices[j].AvailMem = PhysicalDevice.Memory - PhysicalDevice.MemoryUsed
			result[i].Devices[j].AvailDisk = PhysicalDevice.Disk - PhysicalDevice.DiskUsed
			result[i].Devices[j].AvailNetRate = PhysicalDevice.NetRate - PhysicalDevice.NetRateUsed
		}
	}
	coll.RemoveAll(nil)
	return result, nil
}
