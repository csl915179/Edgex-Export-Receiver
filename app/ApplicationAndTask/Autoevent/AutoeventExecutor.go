package Autoevent

import (
	"Edgex-Export_Receiver/app/EdgexData"
	"Edgex-Export_Receiver/app/domain"
	"log"
	"time"
)


type Executor interface {
	Run()
	Stop()
}

type executor struct {
	app				domain.Application
	stop			bool
}

//执行
func (e *executor) Run() {
	for {
		if e.stop {
			break
		}
		if e.app.Frequency == 0 {
			continue
		}
		time.Sleep(time.Duration(e.app.Frequency) * time.Second)
		err := EdgexData.InitEvent(e.app)
		if err != nil {
			log.Println(err)
		}
	}
}

//停止
func (e *executor) Stop() {
	e.stop = true
}

//初始化
func Newexecutor (application domain.Application) Executor {
	return &executor{app:application, stop:(!application.AutoEventState)}
}