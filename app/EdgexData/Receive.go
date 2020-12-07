package EdgexData

import (
	"Edgex-Export_Receiver/app/TaskSchedule"
	"Edgex-Export_Receiver/app/config"
	"Edgex-Export_Receiver/app/db"
	"Edgex-Export_Receiver/app/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)


//处理收到的event,发送调度并放入Mongo
func Receive(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	var application domain.Application
	if err := json.NewDecoder(r.Body).Decode(&application); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	eventElement := domain.Event{}
	eventElement.TranslateApplicationtoEvent(application)
	if _,err := db.GetEventRepos().Insert(&eventElement); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	PostEvent(eventElement)
}

//把event以POST方式发出去
func PostEvent(event domain.Event)  {
	url := "http://" + config.ScheduleConf.AppSchedule.Host + ":" + strconv.FormatInt(config.ScheduleConf.AppSchedule.Port, 10) + "/" + config.ScheduleConf.AppSchedule.Path
	fmt.Println(url)
	client := &http.Client{Timeout: 5*time.Second}
	jsonStr, _ := json.Marshal(event)
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

//从mongo里抽出所有的event,存入待执行列表并返回
func Pull(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	eventlist,err := db.GetEventRepos().ExtractAll()
	if err != nil{
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	result, _ := json.Marshal(eventlist)
	w.Write(result)
	err = db.GetEventToExecuteRepos().InsertIntoToExecute(eventlist)
	if err != nil {
		log.Println(err.Error())
	}
	go TaskSchedule.ReceiveTaskScheduleResult();
}