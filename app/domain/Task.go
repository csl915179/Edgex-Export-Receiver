package domain

import (
	"gopkg.in/mgo.v2/bson"
)

//描述属性用
type Attribute struct {
	//用于描述某个属性
	Key				string				`json:"key"`				//属性名
	Value			string				`json:"value"`				//属性取值
}
type Attribute_Rule struct {
	//用于描述某个属性需要满足的条件
	Key				string				`json:"key"`				//某个key
	Operator		string				`json:"value"`				//某个operator，取值范围:{In, NotIn, Exist, NotExist, GT, LT}
	Value_Num		int64				`json:"value_num"`			//key的取值数目
	Value_List		[]string			`json:"value_list"`			//key的取值列表
}
//对于一个实体（Node或Pod），它可能在某些时候需要罗列自己所有的Attribute或Attribute_Rule
type 	Entity_Attribute	struct{
	Name			string				`json:"name"`				//实体名称
	Attribute		Attribute			`json:"attribute"`			//实体属性
}
type Entity_Attribute_Rule struct {
	Name			string				`json:"name"`				//实体名称
	Attribute_Rule	Attribute_Rule		`json:"attribute_rule"`		//实体属性规则
}

//Task类型，记录某条设备指令对应的资源消耗等情况，收到event后找出对应关系，和device的数据拼接后发走。
type Task struct {
	Id          	bson.ObjectId		`bson:"_id,omitempty" json:"id"`					//Task本身在数据库里的ID
	DeviceId		string				`json:"deviceid"`									//Task对应的设备的ID
	DeviceCommand	string				`json:"devicecommand"`								//Task对应的设备的命令
	Description 	string        		`json:"desc"`										//描述
	HostName		string				`json:"host_name"`									//Task要占用的主机名
	HostPort		string				`json:"host_port"`									//Task要占用的端口号，范围5001-65535
	Kind			string				`json:"kind"`										//Task要建立的Pod类型，取值{"ReplicationController","ReplicaSet","Pod","Service","DaemonSet","Deployment"}，Pod为主
	Tolerations		[]Attribute			`json:"tolerations"`								//Task的容忍列表
	ImageNeed		[]string			`json:"image_need"`									//Task需要的镜像
	CPURequest		int64				`json:"cpu_request"`								//Task需要的CPU资源
	MemoryRequest	int64				`json:"memory_request"`								//Task需要的内存资源
	DiskRequest		int64				`json:"disk_request"`								//Task需要的磁盘资源
	TaskLabels		[]Attribute			`json:"task_labels"`								//Task(Pod)的标签
	ExecLimit   	string       		`json:"exec_limit"`									//执行地点限制 Local/Remote/LocalOrRemote
	TimeLimit		string				`json:"time_limit"`									//完成时间限制，格式为数字+ms/s/min/h/d
}