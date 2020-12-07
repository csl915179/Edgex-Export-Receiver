package Mongo

import (
	"Edgex-Export_Receiver/app/domain"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type EventToExecuteMongoRepository struct {
}

//将Event加入待执行列表
func (ar *EventToExecuteMongoRepository) InsertIntoToExecute (eventList []domain.Event) error {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(eventtoexecuteScheme)
	for i:=0; i<len(eventList); i++ {
		//eventList[i].ScheduleTime = time.Now().Format("2006-01-02T15:04:05")
		err := coll.Insert(eventList[i])
		if err != nil {
			log.Println("Insert event into EventToExecuteCollection Failed! ", err.Error())
		}
	}
	return nil
}

//将Event从待执行列表中抽取出来
func (ar *EventToExecuteMongoRepository) Extract (id string) (domain.Event, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(eventtoexecuteScheme)
	result := domain.Event{}
	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		log.Println("Find event when extract ino EventExecutedCollection failed !", err.Error())
		return result, err
	}
	err = coll.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		log.Println("Remove event when extract ino EventExecutedCollection failed !", err.Error())
		return result, err
	}
	return result, nil
}

//按ID查询Event
func (ar *EventToExecuteMongoRepository) Select (id string) (domain.Event, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(eventtoexecuteScheme)
	result := domain.Event{}
	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		//log.Println("Select event failed !", err.Error())
		return result, err
	}
	return result, nil
}

//查询所有Event
func (ar *EventToExecuteMongoRepository) SelectAll() ([]domain.Event, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(eventtoexecuteScheme)
	result := make([]domain.Event, 0)
	err := coll.Find(nil).Sort("schedule_time").All(&result)
	if err != nil {
		log.Println("Find All EventToExecute failed !" + err.Error())
		return result, err
	}
	return result, err
}

//查询最近的几条Event
func (ar *EventToExecuteMongoRepository) SelectNumber(number int64) ([]domain.Event, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(eventtoexecuteScheme)
	result := make([]domain.Event, 0)
	err := coll.Find(nil).All(&result)
	if err != nil {
		log.Println("Find EventToExecute failed !" + err.Error())
		return result, err
	}
	return result, err
}

//更新Event
func (ar *EventToExecuteMongoRepository) Update (event *domain.Event) (domain.Event, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(eventtoexecuteScheme)
	err := coll.UpdateId(event.Id, &event)
	if err != nil {
		log.Println("Update Device failed !" + err.Error())
		return *event, err
	}
	return *event,nil
}

//删除event
func (ar *EventToExecuteMongoRepository) Delete (id string) error {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(eventtoexecuteScheme)
	err := coll.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		log.Println("Remove event in EventExecutedCollection failed !", err.Error())
		return err
	}
	return nil
}