package Mongo

import (
	"Edgex-Export_Receiver/app/domain"
	"log"
)

type EventExecutedMongoRepository struct {
}

func (ar *EventExecutedMongoRepository) InsertIntoExecuted (event *domain.Event) error {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(eventexecutedSchene)
	err := coll.Insert(event)
	if err != nil {
		//log.Println("Insert event into EventExecuted failed !")
		return err
	}
	return nil
}

//查询所有Event
func (ar *EventExecutedMongoRepository) SelectAll() ([]domain.Event, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(eventtoexecuteScheme)
	result := make([]domain.Event, 0)
	err := coll.Find(nil).Sort("executetime").All(&result)
	if err != nil {
		log.Println("Find All EventToExecute failed !" + err.Error())
		return result, err
	}
	return result, err
}

//查询最新的某几条Event
func (ar *EventExecutedMongoRepository) SelectNumber(number int64) ([]domain.Event, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(eventexecutedSchene)
	result := make([]domain.Event, 0)
	err := coll.Find(nil).Sort("execute_time").All(&result)
	if err != nil {
		log.Println("Find Event Executed failed !" + err.Error())
		return result, err
	}
	return result, err
}