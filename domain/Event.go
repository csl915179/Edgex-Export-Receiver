package domain

import (
	contract	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"gopkg.in/mgo.v2/bson"
)

type Event struct {
	Id          	bson.ObjectId					`bson:"_id,omitempty" json:"id"`
	Event 			contract.Event					`json:"event"`
}