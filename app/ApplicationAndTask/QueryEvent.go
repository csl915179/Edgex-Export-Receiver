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
	low,_ := strconv.Atoi(vars["low"])
	high,_ := strconv.Atoi(vars["high"])

	event,err := db.GetEventRepos().SelectNumber(low, high)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	result, _ := json.Marshal(&event)
	w.Write(result)
}

//查询所有等待调度的Event
func FindEvent (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	event,err := db.GetEventRepos().SelectAll()
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
	low,_ := strconv.Atoi(vars["low"])
	high,_ := strconv.Atoi(vars["high"])

	event,err := db.GetEventToExecuteRepos().SelectNumber(low, high)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	result, _ := json.Marshal(&event)
	w.Write(result)
}

//查询所有返回了调度结果的Event
func FindEventToExecute (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	event,err := db.GetEventToExecuteRepos().SelectAll()
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
	low,_ := strconv.Atoi(vars["low"])
	high,_ := strconv.Atoi(vars["high"])

	event,err := db.GetEventExecutedRepos().SelectNumber(low,high)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	result, _ := json.Marshal(&event)
	w.Write(result)
}

//查询被执行的Event
func FindEventExecuted (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	event,err := db.GetEventExecutedRepos().SelectAll()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	result, _ := json.Marshal(&event)
	w.Write(result)
}