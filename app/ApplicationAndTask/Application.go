package ApplicationAndTask

import (
	"Edgex-Export_Receiver/app/db"
	"Edgex-Export_Receiver/app/domain"
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

//查询所有Application
func ListAllApplication (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	//vars := mux.Vars(r)
	//number := vars["number"]

	ApplicationList,err := db.GetApplicationRepos().SelectAll()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	result, _ := json.Marshal(&ApplicationList)
	w.Write(result)
}

//按ID查询某个Application
func FindApplicationByID (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	application,err := db.GetApplicationRepos().Select(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	result, _ := json.Marshal(&application)
	w.Write(result)
}

//按ID删除某个Application
func DeleteApplicationByID (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	err := db.GetApplicationRepos().Delete(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
}

//增加某个Application
func AddApplication(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	var Application domain.Application
	if err := json.NewDecoder(r.Body).Decode(&Application); err != nil {
		log.Println("Decode Error", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	if Application.Id.Hex() == ""{
		Application.Id = bson.NewObjectId()
	}

	if _,err := db.GetApplicationRepos().Insert(&Application); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	Applicationbody,_ := json.Marshal(Application)
	log.Println(string(Applicationbody))
	w.Write([]byte(Application.Id.Hex()))
}

//修改某个Application
func EditApplication(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	var Application domain.Application
	if err := json.NewDecoder(r.Body).Decode(&Application); err != nil {
		log.Println("Decode Error", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	if _,err := db.GetApplicationRepos().Update(&Application); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	Applicationbody,_ := json.Marshal(Application)
	w.Write(Applicationbody)
}