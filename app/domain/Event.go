package domain

import (
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"time"
)

type EventTaskExecResult string
const (
	OK 		EventTaskExecResult = "OK"
	Fail 	EventTaskExecResult = "Fail"
)

type Event struct {
	Id          	bson.ObjectId						`bson:"_id,omitempty" json:"id"`
	AppID			string								`json:"app_id"`
	Type			string								`json:"type"`
	Name        	string        						`json:"name"`
	Frequency		int64								`json:"frequency"`
	Description 	string        						`json:"desc"`
	Modified		time.Time							`json:"modified"`
	Devices			[]device							`json:"devices"`
}

type device struct {
	Id				string							`json:"id"`
	Name			string							`json:"name"`
	AvailCpu		int64							`json:"avail_cpu"`
	AvailMem		int64							`json:"avail_mem"`
	AvailDisk		int64							`json:"avail_disk"`
	AvailNetRate	int64							`json:"avail_net_rate"`
	Tasks			[]devicetask					`json:"tasks"`
}

type devicetask struct {
	Id					string					`json:"id"`
	Name				string					`json:"name"`								//Task的的名称
	Type				string					`json:"type"`								//GetOrPut
	CPURequest			int64					`json:"cpu"`								//Command需要的CPU资源
	MemoryRequest		int64					`json:"memory"`								//Command需要的内存资源
	DiskRequest			int64					`json:"disk"`								//Command需要的磁盘资源
	Size				int64					`json:"size"`
	NetRate				int64					`json:"net_rate"`
	TaskLabels			[]Attribute				`json:"task_labels"`						//Task(Pod)的标签
	ExecLimit   		string       			`json:"exec_limit"`							//执行地点限制 Local/Remote/LocalOrRemote
	TimeLimit			int64					`json:"time_limit"`							//完成时间限制，单位为S
	EnergyLimit   		int64       			`json:"energy_limit"`						//能耗限制
	ExecPlace			int64					`json:"exec_place"`							//最后实际执行地点
	ExecTime			time.Time				`json:"exec_time"`							//最后实际执行时间
	EnergyUsed 			int64					`json:"energy_used"`						//最后能耗
	ExecResult 			EventTaskExecResult		`json:"exec_result"`						//执行结果
}

//把收到的application转换成Event
func (event *Event) TranslateApplicationtoEvent(application Application) {
	event.Id = bson.NewObjectId()
	event.AppID = application.Id.Hex()
	event.Type = application.Type
	event.Name = application.Name
	event.Frequency = application.Frequency
	event.Description = application.Description
	event.Modified = time.Now()
	event.Devices = make([]device, 0)
	for _,task := range application.DeviceTasks{
		device := device{Id:task.DeviceId, Name:task.DeviceName}
		device.Tasks = make([]devicetask, 0)
		for _,t := range task.Tasks {
			devicetask := devicetask{}
			devicetask.Id = bson.NewObjectId().Hex()
			devicetask.Name = t.Name
			commandval := reflect.ValueOf(&t.Command).Elem()
			commandtype := commandval.Type()
			devicetaskval := reflect.ValueOf(&devicetask).Elem()
			for i:=0; i< commandval.NumField(); i++ {
				name := commandtype.Field(i).Name
				if ok := devicetaskval.FieldByName(name).IsValid(); ok{
					devicetaskval.FieldByName(name).Set(reflect.ValueOf(commandval.Field(i).Interface()))
				}
			}
			device.Tasks = append(device.Tasks, devicetask)
		}
		event.Devices = append(event.Devices, device)
	}
}