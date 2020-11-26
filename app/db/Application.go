package db

import (
	"Edgex-Export_Receiver/app/db/Mongo"
	"Edgex-Export_Receiver/app/domain"
)

type ApplicationRepos interface {
	SelectAll() ([]domain.Application, error)
	Insert(application *domain.Application) (string, error)
	Select(id string) (domain.Application, error)
	Update(application *domain.Application) (domain.Application, error)
	Delete(id string) error
}

func GetApplicationRepos() ApplicationRepos {
	if Mongo.DS.S != nil {
		return ApplicationRepos(&Mongo.ApplicationMongoRepository{})
	} else{
		return nil
	}
}