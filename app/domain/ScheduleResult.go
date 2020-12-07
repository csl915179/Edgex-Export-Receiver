package domain

import (
	"gopkg.in/mgo.v2/bson"
)


type ScheduleResult struct {
	Id          			bson.ObjectId				`bson:"_id,omitempty" json:"id"`
	AppId					string						`json:"app_id"`
	Name					string						`json:"name"`
	Tasks					[]scheduleResultTask		`json:"tasks"`
	ScheduledTime			int64						`json:"scheduled_time"`
}

type scheduleResultTask struct {
	Id						string						`json:"id"`
	Size					int64						`json:"size"`
	ExecLoca				string						`json:"exec_loca"`
	EvalTime				int64						`json:"eval_time"`
	EvalEnergy				int64						`json:"eval_energy"`
	SrcAddr					string						`json:"src_addr"`
	Dst_Addr				string						`json:"dst_addr"`
}