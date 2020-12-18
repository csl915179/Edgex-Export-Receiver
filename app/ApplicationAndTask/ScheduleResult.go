package ApplicationAndTask

import (
	"Edgex-Export_Receiver/app/ApplicationAndTask/Execute"
	"Edgex-Export_Receiver/app/db"
	"Edgex-Export_Receiver/app/domain"
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

func ListScheduleResult (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	scheduleresult,err := db.GetScheduleResultRepos().SelectAll()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	result, _ := json.Marshal(&scheduleresult)
	w.Write(result)
}

func ListScheduleResultByNumber (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	low,_ := strconv.Atoi(vars["low"])
	high,_ := strconv.Atoi(vars["high"])

	event,err := db.GetScheduleResultRepos().SelectNumber(low,high)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	result, _ := json.Marshal(&event)
	w.Write(result)
}

func ReceiveScheduleResult (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	//接收结果与格式转换
	//接收到的类型
	type receivedScheduleResult struct {
		Id          			bson.ObjectId					`bson:"_id,omitempty" json:"id"`
		AppId					string							`json:"app_id"`
		Name					string							`json:"name"`
		Time					string							`json:"time"`
		Tasks					[]domain.ScheduleResultTask		`json:"tasks"`
	}
	var ReceivedScheduleResult receivedScheduleResult
	if err := json.NewDecoder(r.Body).Decode(&ReceivedScheduleResult); err != nil {
		log.Println("Decode Error", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	//格式转换
	ScheduleResult := domain.ScheduleResult{}
	ReceivedVal := reflect.ValueOf(&ReceivedScheduleResult).Elem()
	ReceivedType := ReceivedVal.Type()
	ScheduleResultVal := reflect.ValueOf(&ScheduleResult).Elem()
	for i:=0; i<ReceivedVal.NumField(); i++ {
		name := ReceivedType.Field(i).Name
		if ok := ScheduleResultVal.FieldByName(name).IsValid(); ok{
			ScheduleResultVal.FieldByName(name).Set(reflect.ValueOf(ReceivedVal.Field(i).Interface()))
		}
	}
	//时间解析
	ScheduleResult.ShceduledTime, _ = time.ParseInLocation("2006-01-02T15:04:05.000000", ReceivedScheduleResult.Time, time.Local)

	db.GetScheduleResultRepos().Insert(&ScheduleResult)
	db.GetEventRepos().Extract(ScheduleResult.Id.Hex())
	go Execute.GetEventExecutemanager().ExecuteEvent(ScheduleResult.Id.Hex())
	w.Write([]byte("OK"))

	log.Println("Get Schedule Result: ",ScheduleResult)
}