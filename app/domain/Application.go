package domain

import "gopkg.in/mgo.v2/bson"

type Application struct {
	Id          	bson.ObjectId 						`bson:"_id,omitempty" json:"id"`
	Type			string								`json:"type"`
	Name        	string        						`json:"name"`
	Frequency		int64								`json:"frequency"`
	Description 	string        						`json:"desc"`
	Tasks			map[string]map[string]Command		`json:"tasks"`
}