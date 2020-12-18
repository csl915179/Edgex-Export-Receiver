package Mongo

import (
	"Edgex-Export_Receiver/app/domain"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type EventExecutedMongoRepository struct {
}

func (ar *EventExecutedMongoRepository) InsertIntoExecuted (event *domain.Event) error {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(eventexecutedSchene)
	coll.Remove(bson.M{"_id": event.Id})
	coll = ds.S.DB(database).C(eventexecutedSchene)
	event.Modified = time.Now()
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
func (ar *EventExecutedMongoRepository) SelectNumber(low, high int) ([]domain.Event, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(eventexecutedSchene)
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