package EdgexData

import (
	"Edgex-Export_Receiver/app/config"
	"Edgex-Export_Receiver/app/db"
	"Edgex-Export_Receiver/app/domain"
	"bytes"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"net/http"
	"reflect"
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
	eventElement, err := TranslateApplicationtoEvent(application)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if _, err = db.GetEventRepos().Insert(&eventElement); err != nil {
		log.Println(err)
		return err
	}
	if err := PostEvent(eventElement); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func TranslateApplicationtoEvent(application domain.Application) (domain.Event, error) {
	event := domain.Event{}
	event.Id = bson.NewObjectId()
	event.AppID = application.Id.Hex()
	event.Type = application.Type
	event.Name = application.Name
	event.Frequency = application.Frequency
	event.Description = application.Description
	event.Modified = time.Now()
	event.Devices = make([]domain.Eventdevice, 0)
	count := int32(0)
	ran := int32(2) + rand.Int31n(int32(len(application.DeviceTasks)))
	for _, devicetask := range application.DeviceTasks {
		deviceInfo, err := db.GetDeviceRepos().Select(devicetask.DeviceId)
		if err != nil {
			log.Println(err.Error())
			return event, err
		}
		device := domain.Eventdevice{Id: devicetask.DeviceId, Name: devicetask.DeviceName, AvailCpu: deviceInfo.Cpu - deviceInfo.CpuUsed,
			AvailMem: deviceInfo.Memory - deviceInfo.MemoryUsed, AvailDisk: deviceInfo.Disk - deviceInfo.DiskUsed, AvailNetRate: deviceInfo.NetRate - deviceInfo.NetRateUsed}
		device.Tasks = make([]domain.Eventdevicetask, 0)
		for _, t := range devicetask.Tasks {
			count++
			if count%ran == 0 {
				continue
			}
			devicetask := domain.Eventdevicetask{}
			devicetask.Id = bson.NewObjectId().Hex()
			commandval := reflect.ValueOf(&t.Command).Elem()
			commandtype := commandval.Type()
			devicetaskval := reflect.ValueOf(&devicetask).Elem()
			for i := 0; i < commandval.NumField(); i++ {
				name := commandtype.Field(i).Name
				if ok := devicetaskval.FieldByName(name).IsValid(); ok {
					devicetaskval.FieldByName(name).Set(reflect.ValueOf(commandval.Field(i).Interface()))
				}
			}
			devicetask.Name = t.Name
			//增加随机数
			devicetask.CPURequest += rand.Int63n((5))
			devicetask.MemoryRequest += rand.Int63n((5))
			devicetask.DiskRequest += rand.Int63n((5))
			devicetask.NetRate += rand.Int63n(5)
			devicetask.Size += rand.Int63n(5)
			devicetask.TimeLimit += rand.Int63n(5)
			devicetask.EnergyLimit += rand.Int63n(5)
			device.Tasks = append(device.Tasks, devicetask)
		}
		event.Devices = append(event.Devices, device)
	}
	return event, nil
}

//把event以POST方式发出去
func PostEvent(event domain.Event) error {
	url := "http://" + config.ScheduleConf.AppSchedule.Host + ":" + strconv.FormatInt(config.ScheduleConf.AppSchedule.Port, 10) + "/" + config.ScheduleConf.AppSchedule.Path
	log.Println("Pushed event", event, "To", url)
	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, _ := json.Marshal(event)
	_, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
