package domain

import (
	"gopkg.in/mgo.v2/bson"
)

//每当有新的event来到的时候，都会对应产生一个TaskEvent，然后发给边缘去决策或者在本地决策后去执行
type TaskEvent struct {
	Id          	bson.ObjectId		`bson:"_id,omitempty" json:"id"`					//Task本身在数据库里的ID
	Command			Command				`json:"Command"`									//Task对应的Command,记录硬件开销等信息
	ExecPlace		string				`json:"exec_place"`									//记录最后执行的地点
	State       	string       		`json:"exec_state"`									//执行状态 NOT EXECUTED/EXECTUTING/EXECUTED
	Data			string				`json:"data"`										//Task携带的data
}