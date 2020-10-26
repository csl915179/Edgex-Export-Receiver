package domain

import (
	"gopkg.in/mgo.v2/bson"
)

type Command struct {
	Id          	bson.ObjectId		`bson:"_id,omitempty" json:"id"`			//Command本身在数据库里的ID
	DeviceId		string				`json:"deviceid"`							//对应的DeviceID,为Device在Exporter里的ID
	URL				string				`json:"url"`								//指令地址
	Name			string				`json:"name"`								//command的名称
	Type			string				`json:"type"`								//GetOrPut
	Description 	string        		`json:"desc"`								//描述
	HostName		string				`json:"host_name"`							//Command要占用的主机名
	HostPort		string				`json:"host_port"`							//Command要占用的端口号，范围5001-65535
	Kind			string				`json:"kind"`								//Command要建立的Pod类型，取值{"ReplicationController","ReplicaSet","Pod","Service","DaemonSet","Deployment"}，Pod为主
	Tolerations		[]Attribute			`json:"tolerations"`						//Command的容忍列表
	ImageNeed		[]string			`json:"image_need"`							//Command需要的镜像
	CPURequest		int64				`json:"cpu_request"`						//Command需要的CPU资源
	MemoryRequest	int64				`json:"memory_request"`						//Command需要的内存资源
	DiskRequest		int64				`json:"disk_request"`						//Command需要的磁盘资源
	TaskLabels		[]Attribute			`json:"task_labels"`						//Task(Pod)的标签
	ExecLimit   	string       		`json:"exec_limit"`							//执行地点限制 Local/Remote/LocalOrRemote
	TimeLimit		string				`json:"time_limit"`							//完成时间限制，格式为数字+ms/s/min/h/d
}