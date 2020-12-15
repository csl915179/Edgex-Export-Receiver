package domain

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)


type ScheduleResult struct {
	Id          			bson.ObjectId				`bson:"_id,omitempty" json:"id"`		//结果的ID，与之前Event的ID对应
	AppId					string						`json:"app_id"`							//对应应用的ID
	Name					string						`json:"name"`							//对应应用的名称
	ShceduledTime			time.Time					`json:"shceduled_time"`					//调度时间
	Tasks					[]ScheduleResultTask		`json:"tasks"`							//每个应用的调度结果
}

type ScheduleResultTask struct {
	Id						string						`json:"id"`								//对应应用的ID
	Name					string						`json:"name"`							//对应应用的名称
	Size					int64						`json:"size"`
	ExecLoca				string						`json:"exec_local"`
	EvalTime				int64						`json:"exec_time"`
	EvalEnergy				int64						`json:"exec_energy"`
	SrcAddr					string						`json:"src_addr"`						//预留，任务来源
	DstAddr					string						`json:"dst_addr"`						//预留，任务去向
}