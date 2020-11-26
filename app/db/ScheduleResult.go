package db

import (
"Edgex-Export_Receiver/app/db/Mongo"
"Edgex-Export_Receiver/app/domain"
)

type ScheduleResultRepos interface {
	Select (id string) (domain.ScheduleResult, error)
	Insert (scheduleResult *domain.ScheduleResult) (string, error)
	Update (scheduleResult *domain.ScheduleResult) (domain.ScheduleResult,error)
	Delete (id string) error
}

func GetScheduleResultRepos() ScheduleResultRepos {
	if Mongo.DS.S != nil {
		return ScheduleResultRepos(&Mongo.ScheduleResultMongoRepository{})
	} else{
		return nil
	}
}
