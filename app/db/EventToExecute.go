package db

import (
	"Edgex-Export_Receiver/app/db/Mongo"
	"Edgex-Export_Receiver/app/domain"
)

type EventToExecuteRepos interface {
	InsertIntoToExecute (event *domain.Event) error
	Extract (id string) (domain.Event, error)
	Select (id string) (domain.Event, error)
	SelectAll() ([]domain.Event, error)
	SelectNumber(low, high int) ([]domain.Event, error)
	Update (event *domain.Event) (domain.Event, error)
	Delete (id string) error
}

func GetEventToExecuteRepos() EventToExecuteRepos {
	if Mongo.DS.S != nil {
		return EventToExecuteRepos(&Mongo.EventToExecuteMongoRepository{})
	} else{
		return nil
	}
}
