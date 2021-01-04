//管理每一条Event的执行
package Execute

import (
	"Edgex-Export_Receiver/app/db"
	"Edgex-Export_Receiver/app/domain"
	"gopkg.in/mgo.v2/bson"
	"sync"
)

//全局变量
type EventExecute struct {
	createOnce sync.Once
	m          *eventExecutemanager
	mutex      sync.Mutex
}

var eventExecute EventExecute

type eventExecuteUnit struct {
	event     domain.Event
	taskunits map[string]*TaskExecuteUnit
}

type eventExecutemanager struct {
	execsMap  map[string]eventExecuteUnit //记录当前正在执行的Event
	startOnce sync.Once
}

type EventExecutemanager interface {
	ExecuteEvent(id string)
}

//执行某条Event
func (m *eventExecutemanager) ExecuteEvent(id string) {
	event, _ := db.GetEventToExecuteRepos().Select(id)
	eventExecute.mutex.Lock()
	eventExecute.m.execsMap[event.Id.Hex()] = eventExecuteUnit{event: event, taskunits: make(map[string]*TaskExecuteUnit)}
	eventExecute.mutex.Unlock()
	app, _ := db.GetApplicationRepos().Select(event.AppID)
	//将Event对应的所有的Task装入并发送给对应设备
	for _, device := range app.DeviceTasks {
		deviceTaskUnits := make([]*TaskExecuteUnit, 0)
		for _, task := range device.Tasks {
			taskunit := &TaskExecuteUnit{Id: bson.NewObjectId(), Task: task, DeviceId: device.DeviceId, State: not_executed}
			eventExecute.mutex.Lock()
			eventExecute.m.execsMap[event.Id.Hex()].taskunits[taskunit.Id.Hex()] = taskunit
			eventExecute.mutex.Unlock()
			deviceTaskUnits = append(deviceTaskUnits, taskunit)
		}
		go GetDeviceManager().ExecuteTasks(device.DeviceId, deviceTaskUnits)
	}
	//监听，等待全部执行完毕，虽然可能会有执行失败的
	flag := false
	taskunitList := eventExecute.m.execsMap[event.Id.Hex()].taskunits
	for flag == false {
		flag = true
		for _, task := range taskunitList {
			if task.State != executed {
				flag = false
				break
			}
		}
	}
	//将执行结果写入Event
	//收集执行结果
	executedTaskMap := make(map[string]map[string]*TaskExecuteUnit, 0)
	for _, task := range eventExecute.m.execsMap[event.Id.Hex()].taskunits {
		if _, ok := executedTaskMap[task.DeviceId]; !ok {
			executedTaskMap[task.DeviceId] = make(map[string]*TaskExecuteUnit, 0)
		}
		executedTaskMap[task.DeviceId][task.Task.Name] = task
	}
	//把执行结果写入Event
	for i := 0; i < len(event.Devices); i++ {
		for j := 0; j < len(event.Devices[i].Tasks); j++ {
			event.Devices[i].Tasks[j].EnergyUsed = executedTaskMap[event.Devices[i].Id][event.Devices[i].Tasks[j].Name].Result.EnergyUsed
			event.Devices[i].Tasks[j].ExecPlace = executedTaskMap[event.Devices[i].Id][event.Devices[i].Tasks[j].Name].Result.ExecPlace
			event.Devices[i].Tasks[j].ExecTime = executedTaskMap[event.Devices[i].Id][event.Devices[i].Tasks[j].Name].Result.ExecTime
			event.Devices[i].Tasks[j].ExecResult = executedTaskMap[event.Devices[i].Id][event.Devices[i].Tasks[j].Name].Result.ExecResult
		}
	}
	//event入库，Manager清理
	db.GetEventExecutedRepos().InsertIntoExecuted(&event)
	eventExecute.mutex.Lock()
	delete(eventExecute.m.execsMap, event.Id.Hex())
	eventExecute.mutex.Unlock()
}

//初始化EventExecutemanager
func GetEventExecutemanager() EventExecutemanager {
	eventExecute.createOnce.Do(func() {
		eventExecute.m = &eventExecutemanager{execsMap: make(map[string]eventExecuteUnit)}
	})
	return eventExecute.m
}
