package domain

import (
	"gopkg.in/mgo.v2/bson"
)

type Device struct {
	Id          	bson.ObjectId		`bson:"_id,omitempty" json:"id"`					//Device本身在数据库里的ID
	EdgexId			string				`json:"edgexid"`									//Device在Edgex里的ID
	Name			string				`json:"name"`										//device的名称
	Commands		[]Command			`json:"commands"`									//device的command
}
