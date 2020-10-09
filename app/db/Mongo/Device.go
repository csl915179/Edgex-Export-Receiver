package Mongo

import (
	"Edgex-Export_Receiver/app/domain"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type DeviceMongoRepository struct {
}

func (ar *DeviceMongoRepository) SelectAll() ([]domain.Device, error){
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(deviceScheme)
	result := make([]domain.Device, 0)
	err := coll.Find(nil).All(&result)
	if err != nil {
		log.Println("Find All Device failed !" + err.Error())
		return result, err
	}
	return result, err
}

func (ar *DeviceMongoRepository) Select(id string) (domain.Device, error){
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(deviceScheme)
	result := domain.Device{}
	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		log.Println("Select Device failed !" + err.Error())
		return result, err
	}
	return result, nil
}

func (ar *DeviceMongoRepository) SelectByName(name string) (domain.Device, error){
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(deviceScheme)
	result := domain.Device{}
	err := coll.Find(bson.M{"name": name}).One(&result)
	if err != nil {
		log.Println("Select Device failed !" + err.Error())
		return result, err
	}
	return result, nil
}

//按EdgexID查找Device
func (ar *DeviceMongoRepository) SelectByEdgexId(edgexId string) (domain.Device, error){
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(deviceScheme)
	result := domain.Device{}
	err := coll.Find(bson.M{"edgexid": edgexId}).One(&result)
	if err != nil {
		log.Println("Select Device failed !" + err.Error())
		return result, err
	}
	return result, nil
}

func (ar *DeviceMongoRepository) Insert(device *domain.Device) (string, error){
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(deviceScheme)
	if device.Id == ""{
		device.Id = bson.NewObjectId()
	}
	count,_ := coll.Find(bson.M{"name": device.Name}).Count()
	if count>0 {
		err := errors.New("duplicate device")
		log.Println("Insert Device failed !" + err.Error())
		return "", err
	}
	err := coll.Insert(device)
	if err != nil {
		log.Println("Insert Device failed !" + err.Error())
		return "", err
	}
	return device.Id.Hex(), nil
}

func (ar *DeviceMongoRepository) Delete(id string) error{
	ds := DS.DataStore()
	defer ds.S.Close()
	//先删掉对应的Command
	coll := ds.S.DB(database).C(commandScheme)
	_,err := coll.RemoveAll(bson.M{"deviceid":id})
	if err != nil {
		log.Println("Delete device commands failed!" + err.Error())
		return err
	}
	//正片开始，删除device
	coll = ds.S.DB(database).C(deviceScheme)
	err = coll.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		log.Println("Delete device failed!" + err.Error())
		return err
	}
	return nil
}

func (ar *DeviceMongoRepository) Update(device *domain.Device) (domain.Device, error){
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(deviceScheme)

	err := coll.UpdateId(device.Id, &device)
	if err != nil {
		log.Println("Insert Device failed !" + err.Error())
		return *device, err
	}
	return *device, nil
}

//按照Edgex里的DeviceID删掉device
func (ar *DeviceMongoRepository) DeleteByEdgexId(edgexId string) error{
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(deviceScheme)
	device := domain.Device{}
	err := coll.Find(bson.M{"edgexid": edgexId}).One(&device)
	if err != nil {
		log.Println("Select Device failed !" + err.Error())
		return err
	}
	//先删掉对应的Command
	coll = ds.S.DB(database).C(commandScheme)
	_,err = coll.RemoveAll(bson.M{"deviceid":device.Id.Hex()})
	if err != nil {
		log.Println("Delete device commands failed!" + err.Error())
		return err
	}
	coll = ds.S.DB(database).C(deviceScheme)
	err = coll.Remove(bson.M{"edgexid": edgexId})
	if err != nil {
		log.Println("Delete device failed!" + err.Error())
		return err
	}
	return nil
}
