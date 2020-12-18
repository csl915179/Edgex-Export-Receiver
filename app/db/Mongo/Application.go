package Mongo

import (
	"Edgex-Export_Receiver/app/domain"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type ApplicationMongoRepository struct {
}

//查找所有Application
func (ar *ApplicationMongoRepository) SelectAll() ([]domain.Application, error){
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(applicationScheme)
	result := make([]domain.Application, 0)
	err := coll.Find(nil).All(&result)
	if err != nil {
		log.Println("Find All Application failed !" + err.Error())
		return result, err
	}
	return result, err
}

//新建Application
func (ar *ApplicationMongoRepository) Insert(application *domain.Application) (string, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(applicationScheme)
	err := coll.Insert(application)
	if err != nil {
		log.Println("Insert application failed !")
		return "", err
	}
	return application.Id.Hex(), nil
}

//按ID查找Application
func (ar *ApplicationMongoRepository) Select(id string) (domain.Application, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(applicationScheme)
	result := domain.Application{}
	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		log.Println("Select application failed !", err.Error())
		return result, err
	}
	return result, nil
}

//更新Application
func (ar *ApplicationMongoRepository) Update(application *domain.Application) (domain.Application, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(applicationScheme)
	err := coll.UpdateId(application.Id, &application)
	if err != nil {
		log.Println("Update application failed !" + err.Error())
		return *application, err
	}
	return *application, nil
}

//删除Application
func (ar *ApplicationMongoRepository) Delete(id string) error {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(applicationScheme)
	err := coll.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		log.Println("Delete application failed!" + err.Error())
		return err
	}
	return nil
}