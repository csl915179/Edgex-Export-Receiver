package db

import (
	"Edgex-Export_Receiver/app/db/Mongo"
	"Edgex-Export_Receiver/app/domain"
)

type EventRepos interface {
	Select(id string) (domain.Event, error)
	SelectAll() ([]domain.Event, error)
	SelectNumber(low,high int) ([]domain.Event, error)
	Insert(event *domain.Event) (string, error)
	Extract(id string) error
}

func GetEventRepos() EventRepos {
	if Mongo.DS.S != nil {
		return EventRepos(&Mongo.EventMongoRepository{})
	} else{
		return nil
	}
}
