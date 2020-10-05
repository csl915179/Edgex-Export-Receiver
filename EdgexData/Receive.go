package EdgexData

import (
	"Edgex-Export_Receiver/db"
	"Edgex-Export_Receiver/domain"
	"encoding/json"
	"fmt"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
	"io/ioutil"
	"log"
	"net/http"
)


//处理收到的event放入mongo
func Receive(w http.ResponseWriter, r *http.Request) {
	str, _ := readBodyAsString(w, r)
	var event contract.Event
	json.Unmarshal([]byte(str), &event)
	//ReceiveList.Append(event)
	eventElement := domain.Event{Event:event}
	if _,err := db.GetEventRepos().Insert(&eventElement); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	log.Println(str)
}
func readBodyAsString(writer http.ResponseWriter, request *http.Request) (string, error) {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return "", err
	}
	if len(body) == 0 {
		return "", fmt.Errorf("no request body provided")
	}
	return string(body), nil
}


//从mongo里抽出所有的event并返回
func Pull(w http.ResponseWriter, request *http.Request){
	defer request.Body.Close()
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