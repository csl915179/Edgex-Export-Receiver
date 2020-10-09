package EdgexData

import (
	"Edgex-Export_Receiver/app/db"
	"Edgex-Export_Receiver/app/domain"
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

//新建一条Command
func AddCommand(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	var Command domain.Command
	if err := json.NewDecoder(r.Body).Decode(&Command); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	if Command.Id == ""{
		Command.Id = bson.NewObjectId()
	}

	if _,err := db.GetCommandRepos().Insert(&Command); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	commandbody,_ := json.Marshal(Command)
	log.Println(string(commandbody))
	w.Write([]byte(Command.Id.Hex()))
}

//列出所有Command
func ListAllCommand(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	CommandList, err := db.GetCommandRepos().SelectAll()
	if err != nil {
		log.Println(err.Error())
		w.Write([]byte(err.Error()))
		return
	}
	result, _ := json.Marshal(&CommandList)
	w.Write(result)
}

//修改Command
func EditCommand(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	var Command domain.Command
	if err := json.NewDecoder(r.Body).Decode(&Command); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	Command,err := db.GetCommandRepos().Update(&Command)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	result,_ := json.Marshal(Command)
	w.Write(result)
}


//删除Command
func DeleteCommand(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()
	vars := mux.Vars(r)
	id := vars["id"]
	err := db.GetCommandRepos().Delete(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "err.Error()", http.StatusServiceUnavailable)
		return
	}
	w.Write([]byte("OK"))
}