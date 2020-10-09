package Mongo

import (
	"Edgex-Export_Receiver/app/domain"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type CommandMongoRepository struct {
}


func (ar *CommandMongoRepository) Insert(Command *domain.Command) (string, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(commandScheme)
	err := coll.Insert(Command)
	if err != nil {
		log.Println("Insert Command failed !")
		return "", err
	}
	return Command.Id.Hex(), nil
}

func (ar *CommandMongoRepository) Select(id string) (domain.Command, error){
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(commandScheme)
	result := domain.Command{}
	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		log.Println("Insert event failed !")
		return result, err
	}
	return result, nil
}

func (ar *CommandMongoRepository) SelectAll() ([]domain.Command, error){
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(commandScheme)

	result := make([]domain.Command, 0)
	err := coll.Find(nil).All(&result)
	if err != nil {
		log.Println("SelectAll failed!")
		return nil, err
	}
	return result, nil
}

func (ar *CommandMongoRepository) Update(Command *domain.Command) (domain.Command, error){
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(commandScheme)
	err := coll.UpdateId(Command.Id, &Command)
	if err != nil {
		log.Println("Update Command failed !" + err.Error())
		return *Command, err
	}
	return *Command, nil
}

func (ar *CommandMongoRepository) Delete(id string) error {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB(database).C(commandScheme)
	err := coll.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		log.Println("Delete Command failed !")
		return err
	}
	return nil
}
