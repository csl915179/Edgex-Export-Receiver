package domain

import "gopkg.in/mgo.v2/bson"

type Application struct {
	Id          	bson.ObjectId 						`bson:"_id,omitempty" json:"id"`
	Type			string								`json:"type"`
	Name        	string        						`json:"name"`
	Frequency		int64								`json:"frequency"`
	Description 	string        						`json:"desc"`
	AutoEventState	bool								`json:"auto_event_state"`	//Autoevent此时是否正常有效
	DeviceTasks		map[string]DeviceTask				`json:"devicetasks"`
}

type DeviceTask struct {
	DeviceName		string								`json:"device_name"`
	DeviceId		string								`json:"device_id"`
	Tasks			map[string]Task						`json:"tasks"`
}

type Task struct {
	Name		string									`json:"name"`
	Command		Command									`json:"command"`
}