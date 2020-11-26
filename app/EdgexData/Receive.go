package EdgexData

import (
	"Edgex-Export_Receiver/app/TaskSchedule"
	"Edgex-Export_Receiver/app/db"
	"Edgex-Export_Receiver/app/domain"
	"encoding/json"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
	"log"
	"net/http"
)


//处理收到的event放入mongo
func Receive(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	var event contract.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	eventElement := domain.Event{Event:event}
	eventElement.Size = 0
	for i := 0; i < len(event.Readings); i++  {
		eventElement.Size += int64(len(event.Readings[i].Value))
	}
	if _,err := db.GetEventRepos().Insert(&eventElement); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
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