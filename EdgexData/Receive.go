package EdgexData

import (
	"encoding/json"
	"fmt"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
	"io/ioutil"
	"log"
	"net/http"
)

var ReceiveList EventList

func Receive(w http.ResponseWriter, r *http.Request) {
	str, _ := readBodyAsString(w, r)
	log.Println(str)
	var event contract.Event
	json.Unmarshal([]byte(str), &event)
	ReceiveList.Append(event)
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

func Pull(w http.ResponseWriter, request *http.Request){
	defer request.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	result, _ := json.Marshal(ReceiveList.Clear())
	w.Write(result)
}