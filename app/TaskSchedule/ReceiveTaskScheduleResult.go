package TaskSchedule

import (
	"Edgex-Export_Receiver/app/config"
	"Edgex-Export_Receiver/app/db"
	"Edgex-Export_Receiver/app/domain"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func ReceiveTaskScheduleResult () error{
	time.Sleep(5 * time.Second)
	url := "http://" + config.ScheduleConf.GetSchedule.Host + ":" + strconv.FormatInt(config.ScheduleConf.GetSchedule.Port, 10) + "/" + config.ScheduleConf.GetSchedule.Path
	res, err :=http.Get(url)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	receive, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	var result []domain.ScheduleResult
	err = json.Unmarshal(receive, &result)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	for i:=0; i<len(result); i++ {
		result[i].AppId = result[i].Id.Hex()
		result[i].Id = bson.NewObjectId()
		result[i].ScheduledTime = time.Now().UnixNano()
		db.GetScheduleResultRepos().Insert(&result[i])
	}
	return nil
}