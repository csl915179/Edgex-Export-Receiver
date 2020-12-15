package Mongo

import (
	"Edgex-Export_Receiver/app/domain"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type EventMongoRepository struct {
}

func (ar *EventMongoRepository) SelectNumber(low,high int) ([]domain.Event, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(eventScheme)
	count,_ := coll.Find(nil).Count()
	if low>high {
		low, high = high, low
	}
	if low<0 {
		low = 0
	}
	if low >= count {
		low = count
	}
	if high >= count {
		high = count
	}
	result := make([]domain.Event, 0)
	err := coll.Find(nil).Sort("-modified").Skip(low).Limit(high-low).All(&result)
	if err != nil {
		log.Println("Find Events failed !" + err.Error())
		return result, err
	}
	return result, nil
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

func (ar *EventMongoRepository) SelectAll() ([]domain.Event, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(eventScheme)
	result := make([]domain.Event, 0)
	err := coll.Find(nil).Sort("-executetime").All(&result)
	if err != nil {
		log.Println("Find All Event failed !" + err.Error())
		return result, err
	}
	return result, err
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

func (ar *EventMongoRepository) Extract(id string) error {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(eventScheme)
	result := domain.Event{}
	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = coll.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		log.Println(err.Error())
		return err
	}
	coll = ds.S.DB(database).C(eventtoexecuteScheme)
	err = coll.Insert(result)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
