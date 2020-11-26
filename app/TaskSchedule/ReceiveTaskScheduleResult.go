package TaskSchedule

import (
	"Edgex-Export_Receiver/app/config"
	"Edgex-Export_Receiver/app/db"
	"Edgex-Export_Receiver/app/domain"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type ScheduleResult struct {
	TaskName				string						`json:"task_name"`
	TaskDescription			string						`json:"task_description"`
	ScheduleResult			string						`json:"schedule_result"`
	ScheduleAlgorithm		string						`json:"schedule_algorithm"`
	ScheduleTime			string						`json:"schedule_time"`
	TaskSource				string						`json:"task_source"`
	CpuRequest				int64						`json:"cpu_request"`
	MemoryRequest			int64						`json:"memory_request"`
	DiskRequest				int64						`json:"disk_request"`

}

//接收调度结果
func ReceiveTaskScheduleResult () error{
	time.Sleep(5 * time.Second)
	url := "http://" + config.ScheduleConf.Host + ":" + strconv.FormatInt(config.ScheduleConf.Port, 10) + "/" + config.ScheduleConf.GetSchedule
	res, err :=http.Get(url)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = ParseTaskScheduleResult(result)
	return nil
}

//把调度结果解析成本地格式
func ParseTaskScheduleResult (getScheduleResult []byte) error{
	var result []ScheduleResult
	err := json.Unmarshal(getScheduleResult, &result)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	for i:=0; i<len(result); i++ {
		ScheduleResultOutput := TransportToDomainScheduleResult(result[i])
		db.GetScheduleResultRepos().Insert(&ScheduleResultOutput)
		go ExecuteEvent(ScheduleResultOutput.Id.Hex())
	}
	return nil
}


func TransportToDomainScheduleResult(ScheduleResultInput ScheduleResult)  domain.ScheduleResult {
	//先定义一个空的domain.ScheduleResultResult
	ScheduleResultOutput := domain.ScheduleResult{}
	ScheduleResultInputVal := reflect.ValueOf(&ScheduleResultInput).Elem()
	ScheduleResultInputType := ScheduleResultInputVal.Type()
	ScheduleResultOutputVal := reflect.ValueOf(&ScheduleResultOutput).Elem()

	for i:=0; i<ScheduleResultInputVal.NumField(); i++ {
		name := ScheduleResultInputType.Field(i).Name
		if ok := ScheduleResultOutputVal.FieldByName(name).IsValid(); ok{
			ScheduleResultOutputVal.FieldByName(name).Set(reflect.ValueOf(ScheduleResultInputVal.Field(i).Interface()))
		}
	}
	ScheduleResultOutput.Id =  bson.ObjectIdHex(ScheduleResultOutput.TaskName)
	return ScheduleResultOutput
}