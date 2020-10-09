package EdgexData

import (
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
	if _,err := db.GetEventRepos().Insert(&eventElement); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	eventbody,_ := json.Marshal(event)
	log.Println(string(eventbody))
}


//从mongo里抽出所有的event并返回
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
}