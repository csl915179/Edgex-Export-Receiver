package domain

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type EventTaskExecResult string

const (
	OK   EventTaskExecResult = "OK"
	Fail EventTaskExecResult = "Fail"
)

type Event struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	AppID       string        `json:"app_id"`
	Type        string        `json:"type"`
	Name        string        `json:"name"`
	Frequency   int64         `json:"frequency"`
	Description string        `json:"desc"`
	Modified    time.Time     `json:"modified"`
	Devices     []Eventdevice `json:"devices"`
}

type Eventdevice struct {
	Id           string            `json:"id"`
	Name         string            `json:"name"`
	AvailCpu     int64             `json:"avail_cpu"`
	AvailMem     int64             `json:"avail_mem"`
	AvailDisk    int64             `json:"avail_disk"`
	AvailNetRate int64             `json:"avail_net_rate"`
	Tasks        []Eventdevicetask `json:"tasks"`
}

type Eventdevicetask struct {
	Id            string              `json:"id"`
	Name          string              `json:"name"`   //Task的的名称
	Type          string              `json:"type"`   //GetOrPut
	CPURequest    int64               `json:"cpu"`    //Command需要的CPU资源
	MemoryRequest int64               `json:"memory"` //Command需要的内存资源
	DiskRequest   int64               `json:"disk"`   //Command需要的磁盘资源
	Size          int64               `json:"size"`
	NetRate       int64               `json:"net_rate"`
	TaskLabels    []Attribute         `json:"task_labels"`  //Task(Pod)的标签
	ExecLimit     string              `json:"exec_limit"`   //执行地点限制 Local/Remote/LocalOrRemote
	TimeLimit     int64               `json:"time_limit"`   //完成时间限制，单位为S
	EnergyLimit   int64               `json:"energy_limit"` //能耗限制
	ExecPlace     string              `json:"exec_place"`   //最后实际执行地点
	ExecTime      time.Time           `json:"exec_time"`    //最后实际执行时间
	EnergyUsed    int64               `json:"energy_used"`  //最后能耗
	ExecResult    EventTaskExecResult `json:"exec_result"`  //执行结果
}
