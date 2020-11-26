package Mongo

import (
	"Edgex-Export_Receiver/app/domain"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type ScheduleResultMongoRepository struct {
}

func (ar *ScheduleResultMongoRepository) Select(id string) (domain.ScheduleResult, error){
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(scheduleResultScheme)
	result := domain.ScheduleResult{}
	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		log.Println("Select ScheduleResult failed !" + err.Error())
		return result, err
	}
	return result, nil
}

func (ar *ScheduleResultMongoRepository) Insert(scheduleResult *domain.ScheduleResult) (string, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(scheduleResultScheme)
	err := coll.Insert(scheduleResult)
	if err != nil {
		//log.Println("Insert ScheduleResult failed !" + err.Error())
		return "", err
	}
	return scheduleResult.Id.Hex(), nil
}

func (ar *ScheduleResultMongoRepository) Update(scheduleResult *domain.ScheduleResult) (domain.ScheduleResult, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(scheduleResultScheme)
	err := coll.UpdateId(scheduleResult.Id, &scheduleResult)
	if err != nil {
		//log.Println("Insert ScheduleResult failed !" + err.Error())
		return *scheduleResult, err
	}
	return *scheduleResult, nil
}

func (ar *ScheduleResultMongoRepository) Delete (id string) error {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(scheduleResultScheme)
	err := coll.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		log.Println("Remove scheduleResult failed !", err.Error())
		return err
	}
	return nil
}