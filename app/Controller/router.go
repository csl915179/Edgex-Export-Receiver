package Controller

import (
	"Edgex-Export_Receiver/app/EdgexData"
	mux "github.com/gorilla/mux"
	"net/http"
)

func InitRestRoutes() http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()

	//收发相关
	s.HandleFunc("/push", EdgexData.Receive).Methods(http.MethodPost)
	s.HandleFunc("/pull", EdgexData.Pull).Methods(http.MethodGet)


	//Device相关
	//列出Device
	s.HandleFunc("/device", EdgexData.ListAllDevice).Methods(http.MethodGet)
	//按名查找Device
	s.HandleFunc("/device/name/{name}", EdgexData.FindDeviceByName).Methods(http.MethodGet)
	//按EdgexID查找Device
	s.HandleFunc("/device/edgexid/{edgexid}", EdgexData.FindDeviceByEdgexId).Methods(http.MethodGet)
	//新建Device
	s.HandleFunc("/device", EdgexData.AddDevice).Methods(http.MethodPost)
	//修改Device
	s.HandleFunc("/device", EdgexData.EditDevice).Methods(http.MethodPut)
	//删除Device
	s.HandleFunc("/device", EdgexData.DeleteDevice).Methods(http.MethodDelete)
	//按EdgexID删除Device
	s.HandleFunc("/deviceedgexid/{edgexid}", EdgexData.DeleteDeviceByEdgexId).Methods(http.MethodDelete)
	//整理僵尸Device并列出结果
	s.HandleFunc("/device/check", EdgexData.CheckDeviceExist).Methods(http.MethodPost)


	//Command相关
	//查找所有Command用get
	s.HandleFunc("/command", EdgexData.ListAllCommand).Methods(http.MethodGet)
	//新建某个Command用post，修改用put
	s.HandleFunc("/command", EdgexData.AddCommand).Methods(http.MethodPost)
	s.HandleFunc("/command", EdgexData.EditCommand).Methods(http.MethodPut)
	//删除某个Command
	s.HandleFunc("/command/{id}", EdgexData.DeleteCommand).Methods(http.MethodDelete)

	return r
}