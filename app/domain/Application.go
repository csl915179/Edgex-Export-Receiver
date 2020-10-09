package domain

import "gopkg.in/mgo.v2/bson"

type Application struct {
	Id          	bson.ObjectId 		`bson:"_id,omitempty" json:"id"`
	Type			string				`json:"type"`										//因为应用必须在适配的Node（手机，电脑，平板等）上执行，所以需要绑定一下具体哪一个本地Node
	Name        	string        		`json:"name"`
	Description 	string        		`json:"desc"`
	Tasks			[]Command			`json:"tasks"`
}