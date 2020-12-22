package Execute

import (
	"Edgex-Export_Receiver/app/db"
	"Edgex-Export_Receiver/app/domain"
	"sort"
	"sync"
	"time"
)

//全局变量
type DeviceManagerVars struct {
	createOnce sync.Once
	m          *devicemanager
	mutex      sync.Mutex
}

var deviceManagerVars DeviceManagerVars

type deviceManage struct {
	lock     *sync.Mutex
	taskList []*TaskExecuteUnit
}
type devicemanager struct {
	deviceMap map[string]deviceManage
}

//从任务列表中拿出某个TaskExecuteUnit
func popTask(index int, deviceid string) *TaskExecuteUnit {
	task := deviceManagerVars.m.deviceMap[deviceid].taskList[index]
	tl := append(deviceManagerVars.m.deviceMap[deviceid].taskList[:index], deviceManagerVars.m.deviceMap[deviceid].taskList[index+1:]...)
	l := deviceManagerVars.m.deviceMap[deviceid].lock
	deviceManagerVars.m.deviceMap[deviceid] = deviceManage{lock: l, taskList: tl}
	return task
}

//从任务列表中抛出执行失败的TaskExecuteUnit
func failTask(index int, deviceid string) {
	task := popTask(index, deviceid)
	task.Result = result{ExecResult: domain.Fail}
	task.State = executed
}

type DeviceManager interface {
	LoadAllDevice()
	ReloadDevices()
	AddDevice(device domain.Device)
	RemoveDevice(id string)
	ExecuteTasks(deviceid string, tasks []*TaskExecuteUnit)
	allocateTasksForDevice(deviceid string)
}

//载入所有设备
func (m *devicemanager) LoadAllDevice() {
	var deviceList, _ = db.GetDeviceRepos().SelectAll()
	deviceManagerVars.mutex.Lock()
	for i := 0; i < len(deviceList); i++ {
		deviceManagerVars.m.deviceMap[deviceList[i].Id.Hex()] = deviceManage{lock: new(sync.Mutex), taskList: make([]*TaskExecuteUnit, 0)}
	}
	deviceManagerVars.mutex.Unlock()
}

//重新载入设备列表
func (m *devicemanager) ReloadDevices() {
	var deviceList, _ = db.GetDeviceRepos().SelectAll()
	deviceManagerVars.mutex.Lock()
	newDeviceMap := make(map[string]deviceManage, 0)
	for i := 0; i < len(deviceList); i++ {
		newDeviceMap[deviceList[i].Id.Hex()] = deviceManagerVars.m.deviceMap[deviceList[i].Id.Hex()]
		for len(deviceManagerVars.m.deviceMap[deviceList[i].Id.Hex()].taskList) > 0 {
			failTask(i, deviceList[0].Id.Hex())
		}
		delete(deviceManagerVars.m.deviceMap, deviceList[i].Id.Hex())
	}
	deviceManagerVars.m.deviceMap = newDeviceMap
	deviceManagerVars.mutex.Unlock()
}

//载入某个设备
func (m *devicemanager) AddDevice(device domain.Device) {
	deviceManagerVars.mutex.Lock()
	if _, ok := deviceManagerVars.m.deviceMap[device.Id.Hex()]; ok {
		deviceManagerVars.mutex.Unlock()
		return
	}
	deviceManagerVars.m.deviceMap[device.Id.Hex()] = deviceManage{lock: new(sync.Mutex), taskList: make([]*TaskExecuteUnit, 0)}
	deviceManagerVars.mutex.Unlock()
}

//删除某个设备
func (m *devicemanager) RemoveDevice(id string) {
	deviceManagerVars.mutex.Lock()
	for len(deviceManagerVars.m.deviceMap[id].taskList) > 0 {
		failTask(0, id)
	}
	delete(deviceManagerVars.m.deviceMap, id)
	deviceManagerVars.mutex.Unlock()
}

//载入某个任务到设备
func (m *devicemanager) ExecuteTasks(deviceid string, tasks []*TaskExecuteUnit) {
	//此时首先锁定设备，因为这时要对设备内部的数组进行修改，还要根据时间限制进行排序，不能冲突
	deviceManagerVars.mutex.Lock()
	if _, ok := deviceManagerVars.m.deviceMap[deviceid]; !ok {
		for i := 0; i < len(tasks); i++ {
			tasks[i].Result = result{ExecResult: domain.Fail}
			tasks[i].State = executed
		}
		deviceManagerVars.mutex.Unlock()
		return
	}
	deviceManagerVars.m.deviceMap[deviceid].lock.Lock()
	newTaskList := append(deviceManagerVars.m.deviceMap[deviceid].taskList, tasks...)
	lock := deviceManagerVars.m.deviceMap[deviceid].lock
	deviceManagerVars.m.deviceMap[deviceid] = deviceManage{lock: lock, taskList: newTaskList}
	m.allocateTasksForDevice(deviceid)
	deviceManagerVars.m.deviceMap[deviceid].lock.Unlock()
	deviceManagerVars.mutex.Unlock()
}

//初始化DeviceManager
func GetDeviceManager() DeviceManager {
	deviceManagerVars.createOnce.Do(func() {
		deviceManagerVars.m = &devicemanager{deviceMap: make(map[string]deviceManage, 0)}
	})
	return deviceManagerVars.m
}

//为某个设备分配任务
func (dm deviceManage) Len() int      { return len(dm.taskList) }
func (dm deviceManage) Swap(i, j int) { dm.taskList[i], dm.taskList[j] = dm.taskList[j], dm.taskList[i] }
func (dm deviceManage) Less(i, j int) bool {
	return dm.taskList[i].Task.Command.TimeLimit < dm.taskList[j].Task.Command.TimeLimit
}
func (m *devicemanager) allocateTasksForDevice(deviceid string) {
	//定义如何执行任务,这里如果条件允许那么只有分配的CPU系数会不同
	execTask := func(device *domain.Device, cpu_count int64, cpu_unit int64, index int) {
		task := popTask(index, deviceid)
		task.Cpu = cpu_unit * cpu_count
		task.Memory = task.Task.Command.MemoryRequest
		task.Disk = task.Task.Command.DiskRequest
		task.NetRate = task.Task.Command.Size
		device.CpuUsed += task.Cpu
		device.MemoryUsed += task.Memory
		device.DiskUsed += task.Disk
		device.NetRateUsed += task.NetRate
		db.GetDeviceRepos().Update(device)
		task.State = executing
		go func() {
			time.Sleep(time.Duration(task.Task.Command.CPURequest/task.Cpu) * time.Second)
			task.Result = result{ExecTime: time.Now(), ExecPlace: "LOCAL", EnergyUsed: task.Task.Command.CPURequest / task.Cpu, ExecResult: domain.OK}
			task.State = executed
			newdevice, _ := db.GetDeviceRepos().Select(device.Id.Hex())
			device = &newdevice
			device.CpuUsed -= task.Cpu
			device.MemoryUsed -= task.Memory
			device.DiskUsed -= task.Disk
			device.NetRateUsed -= task.NetRate
			db.GetDeviceRepos().Update(device)
		}()
	}
	//按照完成时间限制排队,先执行完成时间限制高的,也就是TimeLimit小的
	sort.Sort(deviceManagerVars.m.deviceMap[deviceid])
	//为其分配CPU和内存，先给紧急任务分
	for i := 0; i < len(deviceManagerVars.m.deviceMap[deviceid].taskList); i++ {
		//如果找不到设备，执行失败
		device, err := db.GetDeviceRepos().Select(deviceid)
		if err != nil {
			failTask(i, deviceid)
		}
		task := deviceManagerVars.m.deviceMap[deviceid].taskList[i]
		command := task.Task.Command
		//如果设备的全部资源也不够能完成这些任务，直接抛出执行失败
		if (command.TimeLimit > 0 && command.CPURequest/device.Cpu > command.TimeLimit) ||
			command.DiskRequest > device.Disk || command.MemoryRequest > device.Memory || command.Size > device.NetRate {
			failTask(i, deviceid)
		}
		avail_cpu := device.Cpu - device.CpuUsed
		avail_memory := device.Memory - device.MemoryUsed
		avail_disk := device.Disk - device.DiskUsed
		avail_netrate := device.NetRate - device.NetRateUsed
		//其他条件满足，分配CPU
		if avail_disk >= command.DiskRequest && avail_memory >= command.MemoryRequest && avail_netrate >= command.Size {
			cpu_unit := device.Cpu / 20
			cpu_count := int64(1)
			for ; cpu_count*cpu_unit <= avail_cpu; cpu_count++ {
				if task.Task.Command.TimeLimit == 0 {
					execTask(&device, cpu_count, cpu_unit, i)
					break
				}
				if task.Task.Command.CPURequest/(cpu_count*cpu_unit) < task.Task.Command.TimeLimit {
					execTask(&device, cpu_count, cpu_unit, i)
					break
				}
			}
			task.Cpu = cpu_unit * cpu_count
			task.Memory = task.Task.Command.MemoryRequest
			task.Disk = task.Task.Command.DiskRequest
			task.NetRate = task.Task.Command.Size
		} else {
			continue
		}
	}
}
