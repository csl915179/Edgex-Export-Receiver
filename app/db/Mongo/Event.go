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
		//log.Println("Select event failed !", err.Error())
		return result, err
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
		log.Println("Insert event failed !")
		return result, err
	}
	_,err = coll.RemoveAll(nil)
	if err != nil {
		log.Println("Clear event List failed !")
		return result, err
	}
	//coll.Remove(nil)
	return result, nil
}
