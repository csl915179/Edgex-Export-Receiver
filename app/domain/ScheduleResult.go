package domain

import "gopkg.in/mgo.v2/bson"


type ScheduleResult struct {
	Id          			bson.ObjectId				`bson:"_id,omitempty" json:"id"`
	TaskName				string						`json:"task_name"`
	TaskDescription			string						`json:"task_description"`
	ScheduleResult			string						`json:"schedule_result"`
	ScheduleAlgorithm		string						`json:"schedule_algorithm"`
	ScheduleTime			string						`json:"schedule_time"`
}