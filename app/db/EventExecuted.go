package db

import (
"Edgex-Export_Receiver/app/db/Mongo"
	"Edgex-Export_Receiver/app/domain"
)

type EventExecuted interface {
	InsertIntoExecuted (event *domain.Event) error
	SelectAll() ([]domain.Event, error)
	SelectNumber(low, high int) ([]domain.Event, error)
}

func GetEventExecutedRepos() EventExecuted {
	if Mongo.DS.S != nil {
		return EventExecuted(&Mongo.EventExecutedMongoRepository{})
	} else{
		return nil
	}
}
