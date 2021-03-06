package Controller

import (
	"Edgex-Export_Receiver/app/ApplicationAndTask"
	"Edgex-Export_Receiver/app/EdgexData"
	"Edgex-Export_Receiver/app/config"
	mux "github.com/gorilla/mux"
	"net/http"
)

func InitRestRoutes() http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()

	//PING
	s.HandleFunc("/ping", Ping).Methods(http.MethodGet)

	//收发相关
	s.HandleFunc("/push", EdgexData.Receive).Methods(http.MethodPost)

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
	s.HandleFunc("/device/edgexid/{edgexid}", EdgexData.DeleteDeviceByEdgexId).Methods(http.MethodDelete)
	//整理僵尸Device并列出结果
	s.HandleFunc("/device/check", EdgexData.CheckDeviceExist).Methods(http.MethodPost)
	//按名查看设备已使用的资源
	s.HandleFunc("/device/resourceused/name/{name}", EdgexData.FindDeviceResourceUsedByName).Methods(http.MethodGet)

	//Application相关
	s.HandleFunc("/application", ApplicationAndTask.ListAllApplication).Methods(http.MethodGet)
	s.HandleFunc("/application/{id}", ApplicationAndTask.FindApplicationByID).Methods(http.MethodGet)
	s.HandleFunc("/application/{id}", ApplicationAndTask.DeleteApplicationByID).Methods(http.MethodDelete)
	s.HandleFunc("/application", ApplicationAndTask.AddApplication).Methods(http.MethodPost)
	s.HandleFunc("/application", ApplicationAndTask.EditApplication).Methods(http.MethodPut)

	//ScheduleResult相关
	s.HandleFunc("/scheduleresult", ApplicationAndTask.ListScheduleResult).Methods(http.MethodGet)
	s.HandleFunc("/scheduleresult/{low}/{high}", ApplicationAndTask.ListScheduleResultByNumber).Methods(http.MethodGet)
	s.HandleFunc("/scheduleresult", ApplicationAndTask.ReceiveScheduleResult).Methods(http.MethodPost)
	s.HandleFunc("/scheduleresult/appschedule", config.ListAppSchedule).Methods(http.MethodGet)
	s.HandleFunc("/scheduleresult/appschedule", config.ModifyAppSchedule).Methods(http.MethodPost)

	//Event相关
	//查询最近的几条被发去调度的Event
	s.HandleFunc("/event/event/{low}/{high}", ApplicationAndTask.FindEventByNumber).Methods(http.MethodGet)
	s.HandleFunc("/event/eventtoexecute/{low}/{high}", ApplicationAndTask.FindEventToExecuteByNumber).Methods(http.MethodGet)
	s.HandleFunc("/event/eventexecuted/{low}/{high}", ApplicationAndTask.FindEventExecutedByNumber).Methods(http.MethodGet)
	s.HandleFunc("/event/event", ApplicationAndTask.FindEvent).Methods(http.MethodGet)
	s.HandleFunc("/event/eventtoexecute", ApplicationAndTask.FindEventToExecute).Methods(http.MethodGet)
	s.HandleFunc("/event/eventexecuted", ApplicationAndTask.FindEventExecuted).Methods(http.MethodGet)

	s1 := r.PathPrefix("").Subrouter()
	s1.HandleFunc("/ws", WebSocketHandler)

	return r
}

func Ping(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte("Pong"))
}
