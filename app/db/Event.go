package db

import (
	"Edgex-Export_Receiver/app/db/Mongo"
	"Edgex-Export_Receiver/app/domain"
)

type EventRepos interface {
	Select(id string) (domain.Event, error)
	SelectNumber(number int64) ([]domain.Event, error)
	Insert(event *domain.Event) (string, error)
	ExtractAll()([]domain.Event, error)
}

func GetEventRepos() EventRepos {
	if Mongo.DS.S != nil {
		return EventRepos(&Mongo.EventMongoRepository{})
	} else{
		return nil
	}
}
