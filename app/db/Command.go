package db

import (
	"Edgex-Export_Receiver/app/db/Mongo"
	"Edgex-Export_Receiver/app/domain"
)

type CommandRepos interface {
	Insert(Command *domain.Command) (string, error)
	Select(id string) (domain.Command, error)
	SelectAll() ([]domain.Command, error)
	Update(Command *domain.Command) (domain.Command, error)
	Delete(id string) error
}

func GetCommandRepos() CommandRepos {
	if Mongo.DS.S != nil {
		return CommandRepos(&Mongo.CommandMongoRepository{})
	} else{
		return nil
	}
}

