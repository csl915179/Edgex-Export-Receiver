package domain

import (
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
	"gopkg.in/mgo.v2/bson"
)

type Event struct {
	Id          	bson.ObjectId					`bson:"_id,omitempty" json:"id"`
	Event 			contract.Event					`json:"event"`
	ExecutePlace	string							`json:"execute_place"`
	ExecuteTime		string							`json:"execute_time"`
	ScheduleTime	string							`json:"schedule_time"`
	Size			int64							`json:"size"`
}