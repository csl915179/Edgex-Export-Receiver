package domain

import (
	"gopkg.in/mgo.v2/bson"
)

type Device struct {
	Id          	bson.ObjectId			`bson:"_id,omitempty" json:"id"`				//Device本身在数据库里的ID
	EdgexId			string					`json:"edgexid"`								//Device在Edgex里的ID
	Name			string					`json:"name"`									//device的名称
	Cpu				int64					`json:"cpu"`									//device的CPU
	Memory			int64					`json:"memory"`									//device的内存容量
	Disk			int64					`json:"disk"`									//device的硬盘容量
	NetRate			int64					`json:"net_rate"`								//device的网络速率
	GetCommands		map[string]*Command		`json:"getcommands"`							//device的command,get
	PutCommands		map[string]*Command		`json:"putcommands"`							//device的command,put
	CPUUsed			int64					`json:"cpu_used"`								//已经使用的CPU
	MemoryUsed		int64					`json:"memory_used"`							//已经使用的内存
	DiskUsed		int64					`json:"disk_used"`								//已经使用的硬盘
	NetRateUsed		int64					`json:"net_rate_used"`							//已经使用的网络速率
}

type Command struct {
	URL				string				`json:"url"`								//指令地址
	Name			string				`json:"name"`								//command的名称
	Type			string				`json:"type"`								//GetOrPut
	Description 	string        		`json:"desc"`								//描述
	HostName		string				`json:"host_name"`							//Command要占用的主机名
	HostPort		string				`json:"host_port"`							//Command要占用的端口号，范围5001-65535
	Kind			string				`json:"kind"`								//Command要建立的Pod类型，取值{"ReplicationController","ReplicaSet","Pod","Service","DaemonSet","Deployment"}，Pod为主
	Tolerations		[]Attribute			`json:"tolerations"`						//Command的容忍列表
	ImageNeed		[]string			`json:"image_need"`							//Command需要的镜像
	CPURequest		int64				`json:"cpu"`								//Command需要的CPU资源
	MemoryRequest	int64				`json:"memory"`								//Command需要的内存资源
	DiskRequest		int64				`json:"disk"`								//Command需要的磁盘资源
	Size			int64				`json:"size"`
	TaskLabels		[]Attribute			`json:"task_labels"`						//Task(Pod)的标签
	ExecLimit   	string       		`json:"exec_limit"`							//执行地点限制 Local/Remote/LocalOrRemote
	TimeLimit		int64				`json:"time_limit"`							//完成时间限制，格式为数字+ms/s/min/h/d
	EnergyLimit   	int64       		`json:"energy_limit"`						//执行地点限制 Local/Remote/LocalOrRemote
}