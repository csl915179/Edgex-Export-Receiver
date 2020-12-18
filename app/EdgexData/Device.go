package EdgexData

import (
	"Edgex-Export_Receiver/app/ApplicationAndTask/Execute"
	"Edgex-Export_Receiver/app/db"
	"Edgex-Export_Receiver/app/domain"
	"encoding/json"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

//列出所有的device
func ListAllDevice(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	DeviceList, err := db.GetDeviceRepos().SelectAll()
	if err != nil {
		log.Println(err.Error())
		w.Write([]byte(err.Error()))
		return
	}
	result, _ := json.Marshal(&DeviceList)
	w.Write(result)
}

//新建Device
func AddDevice(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	var Device domain.Device
	if err := json.NewDecoder(r.Body).Decode(&Device); err != nil {
		log.Println("Decode Error", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	if Device.Id.Hex() == ""{
		Device.Id = bson.NewObjectId()
	}

	if _,err := db.GetDeviceRepos().Insert(&Device); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	devicebody,_ := json.Marshal(Device)
	log.Println(string(devicebody))
	Execute.GetDeviceManager().AddDevice(Device)
	w.Write([]byte(Device.Id.Hex()))
}

//修改Device
func EditDevice(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	var Device domain.Device
	if err := json.NewDecoder(r.Body).Decode(&Device); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	Device,err := db.GetDeviceRepos().Update(&Device)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	result,_ := json.Marshal(Device)
	w.Write(result)
}

//删除Device
func DeleteDevice (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	err := db.GetDeviceRepos().Delete(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	Execute.GetDeviceManager().RemoveDevice(id)
	w.Write([]byte("OK"))
}

//按EdgexID删除Device
func DeleteDeviceByEdgexId (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	edgexid := vars["edgexid"]
	err := db.GetDeviceRepos().DeleteByEdgexId(edgexid)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.Write([]byte("OK"))
}

//按名称查找device
func FindDeviceByName (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	name := vars["name"]

	device,err := db.GetDeviceRepos().SelectByName(name)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	result, _ := json.Marshal(&device)
	w.Write(result)
}

//按EdgexID查找Device
func FindDeviceByEdgexId (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	edgexid := vars["edgexid"]

	device,err := db.GetDeviceRepos().SelectByEdgexId(edgexid)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	result, _ := json.Marshal(&device)
	w.Write(result)
}

//清除那些已经在Edgex里面删除但还在mongo里面的僵尸device,然后返回一次当前mongo里面的不僵尸的device
func CheckDeviceExist(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	//找出当前存在Edgex里面的device信息
	var ReceivedDeviceList = make([]contract.Device, 0)
	if err := json.NewDecoder(r.Body).Decode(&ReceivedDeviceList); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	//找出当前Mongo里面所有的device信息
	Existed := make(map[string]bool, 0)
	mongoDeviceList, err := db.GetDeviceRepos().SelectAll()
	if err != nil{
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	for i:=0; i<len(mongoDeviceList); i++ {
		Existed[mongoDeviceList[i].EdgexId] = false
	}
	//比对
	for i:=0; i<len(ReceivedDeviceList); i++ {
		Existed[ReceivedDeviceList[i].Id] = true
	}
	for EdgexdeviceId := range(Existed){
		if Existed[EdgexdeviceId] == false{
			err = db.GetDeviceRepos().DeleteByEdgexId(EdgexdeviceId)
			if err != nil {
				log.Println(err.Error())
				w.Write([]byte(err.Error()))
				return
			}
		}
	}
	DeviceList, err := db.GetDeviceRepos().SelectAll()
	if err != nil {
		log.Println(err.Error())
		w.Write([]byte(err.Error()))
		return
	}
	result, _ := json.Marshal(&DeviceList)
	Execute.GetDeviceManager().ReloadDevices()
	w.Write(result)
}