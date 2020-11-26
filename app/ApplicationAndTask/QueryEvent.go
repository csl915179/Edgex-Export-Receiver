package ApplicationAndTask

import (
	"Edgex-Export_Receiver/app/db"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

//查询最近的几条等待调度的Event
func FindEventByNumber (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	number,_ := strconv.ParseInt(vars["number"], 10, 64)

	event,err := db.GetEventRepos().SelectNumber(number)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	result, _ := json.Marshal(&event)
	w.Write(result)
}

//查询最近的几条被发去调度的Event
func FindEventToExecuteByNumber (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	number,_ := strconv.ParseInt(vars["number"], 10, 64)

	event,err := db.GetEventToExecuteRepos().SelectNumber(number)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	result, _ := json.Marshal(&event)
	w.Write(result)
}

//查询最近的几条被执行的Event
func FindEventExecutedByNumber (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	number,_ := strconv.ParseInt(vars["number"], 10, 64)

	event,err := db.GetEventExecutedRepos().SelectNumber(number)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	result, _ := json.Marshal(&event)
	w.Write(result)
}