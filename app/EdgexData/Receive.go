package EdgexData

import (
	"Edgex-Export_Receiver/app/config"
	"Edgex-Export_Receiver/app/db"
	"Edgex-Export_Receiver/app/domain"
	"bytes"
	"encoding/json"
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
	if err := InitEvent(application); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
}

//把应用转换成Event
func InitEvent(application domain.Application) error {
	eventElement := domain.Event{}
	eventElement.TranslateApplicationtoEvent(application)
	if _,err := db.GetEventRepos().Insert(&eventElement); err != nil {
		log.Println(err)
		return err
	}
	if err := PostEvent(eventElement); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//把event以POST方式发出去
func PostEvent(event domain.Event) error {
	url := "http://" + config.ScheduleConf.AppSchedule.Host + ":" + strconv.FormatInt(config.ScheduleConf.AppSchedule.Port, 10) + "/" + config.ScheduleConf.AppSchedule.Path
	log.Println("Pushed event", event, "To", url)
	client := &http.Client{Timeout: 5*time.Second}
	jsonStr, _ := json.Marshal(event)
	_, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}