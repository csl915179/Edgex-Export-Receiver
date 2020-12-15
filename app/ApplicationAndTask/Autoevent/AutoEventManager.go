package Autoevent

import (
	"Edgex-Export_Receiver/app/db"
	"Edgex-Export_Receiver/app/domain"
	"sync"
)

type Manager interface {
	StartAutoEvents()
	StopAutoEvents()
	ReStartForApp(application domain.Application)
	StopForApp(id string)
	StartForApp(application domain.Application)
}

var (
	createOnce sync.Once
	m          *manager
	mutex      sync.Mutex
)

type manager struct {
	execsMap  map[string]Executor
	startOnce sync.Once
}

//启动全部Autoevent
func (m *manager) StartAutoEvents() {
	mutex.Lock()
	m.startOnce.Do(func() {
		apps,_ := db.GetApplicationRepos().SelectAll()
		for i:=0; i<len(apps); i++ {
			exec := Newexecutor(apps[i])
			go exec.Run()
			m.execsMap[apps[i].Id.Hex()] = exec
		}
	})
	mutex.Unlock()
}

//停止全部Autoevent
func (m *manager) StopAutoEvents() {
	mutex.Lock()
	for k, v := range m.execsMap {
		v.Stop()
		delete(m.execsMap, k)
	}
	mutex.Unlock()
}

//停止某个App的Autoevent
func (m *manager) StopForApp(id string) {
	mutex.Lock()
	execs, ok := m.execsMap[id]
	if ok {
		execs.Stop()
		delete(m.execsMap, id)
	}
	mutex.Unlock()
}

//重新启动某个App的Autoevent
func (m *manager) ReStartForApp(application domain.Application) {
	m.StopForApp(application.Id.Hex())
	mutex.Lock()
	exec := Newexecutor(application)
	go exec.Run()
	m.execsMap[application.Id.Hex()] = exec
	mutex.Unlock()
}

//为某个APP增加Autoevent
func (m *manager) StartForApp(application domain.Application){
	mutex.Lock()
	exec := Newexecutor(application)
	go exec.Run()
	m.execsMap[application.Id.Hex()] = exec
	mutex.Unlock()
}


//初始化Manager
func GetManager() Manager {
	createOnce.Do(func() {
		m = &manager{execsMap: make(map[string]Executor)}
	})
	return m
}