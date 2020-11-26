package TaskSchedule

import (
	"Edgex-Export_Receiver/app/db"
	"Edgex-Export_Receiver/app/domain"
	"fmt"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
	"log"
	"regexp"
	"strconv"
	"time"
)

//锁map，记录每一个设备的互斥管道
type deviceLocks struct {
	LockMap map[string]chan int
}
var Locks deviceLocks

//分拣事件，不是本地执行的直接入库，本地执行的执行
func ExecuteEvent (id string) error{
	ScheduleResultInfo, err := db.GetScheduleResultRepos().Select(id)
	if err != nil{
		log.Println("Find Scheduled Result Failed! " + err.Error())
		return err
	}
	event, err := db.GetEventToExecuteRepos().Select(ScheduleResultInfo.TaskName)
	if err != nil{
		//log.Println("Find Scheduled Event Failed! " + err.Error())
		return err
	}
	if ScheduleResultInfo.ScheduleResult == event.Event.ID {
		//ExecuteLocally(&event)
	}else {
		event.ExecutePlace = ScheduleResultInfo.ScheduleResult
		event.ExecuteTime = ScheduleResultInfo.ScheduleTime
		err = db.GetEventExecutedRepos().InsertIntoExecuted(&event)
		if err != nil{
			//log.Println("Insert event into EventExecuted failed ! " + err.Error())
			return err
		}
		db.GetEventToExecuteRepos().Delete(event.Id.Hex())
	}
	return nil
}

func ExecuteLocally (event *domain.Event) error {
	deviceName := event.Event.Device
	DeviceToUse,err := db.GetDeviceRepos().SelectByName(deviceName)
	if err != nil {
		log.Println("Find the device to use failed! " + err.Error());
		return err
	}
	//检查设备锁是否存在
	if _, ok := Locks.LockMap[deviceName]; !ok {
		Locks.LockMap[deviceName] = make(chan int, 1)
		Locks.LockMap[deviceName] <- 1
	}
	for i:=0; i<len(event.Event.Readings); i++ {
		go ExecuteLocalCommand(event.Event.Readings[i], DeviceToUse, *DeviceToUse.GetCommands[event.Event.Readings[i].Name])
	}

	return nil
}

func ExecuteLocalCommand (reading contract.Reading, deviceToUse domain.Device, deviceCommand domain.Command) error{
	time_limit := parseTime(deviceCommand.TimeLimit)
	var time_left int64
	//算时间余量，分配资源
	//没指出时间限制，分配5%CPU
	if time_limit <= 0 {
		time_left = 9223372036854775807
	}else {
		time_left = time_limit - (time.Now().UnixNano()-reading.Origin)/1000000
	}
	fmt.Println(deviceCommand.Name, deviceToUse.Name, time_left)
	return nil
}
func parseTime(time string) int64{
	timeRE := regexp.MustCompile(`^(\d*)+(ms|s|h|min|d)$`)
	if len(timeRE.FindStringSubmatch(time)) == 0 {
		return -1
	}
	time_value, _ := strconv.ParseInt(timeRE.FindStringSubmatch(time)[1],10,64)
	time_unit := timeRE.FindStringSubmatch(time)[2]
	var result int64
	switch time_unit {
		case "ms":
			result = time_value
			break
		case "s":
			result = time_value*1000
			break
		case "min":
			result = time_value*1000*60
			break
		case "h":
			result = time_value*1000*60*60
			break
		case "d":
			result = time_value*1000*60*60*24
			break
	}
	return result
}