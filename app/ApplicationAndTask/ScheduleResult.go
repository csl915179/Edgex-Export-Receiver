package ApplicationAndTask

import (
	"Edgex-Export_Receiver/app/db"
	"Edgex-Export_Receiver/app/domain"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
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
	number,_ := strconv.ParseInt(vars["number"], 10, 64)

	event,err := db.GetScheduleResultRepos().SelectNumber(number)
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

	var Application domain.Application
	if err := json.NewDecoder(r.Body).Decode(&Application); err != nil {
		log.Println("Decode Error", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	log.Println("PONG")
	log.Println(Application)
}