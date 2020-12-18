//管理每一条Task的执行
package Execute

import (
	"Edgex-Export_Receiver/app/domain"
	"gopkg.in/mgo.v2/bson"
	"time"
)

//描述任务执行状态
type executeState int
const (
	not_executed					executeState = 0			//尚未开始执行
	executing						executeState = 1			//执行中
	executed						executeState = 2			//成功执行完成
)

type result struct {
	ExecPlace  int64
	ExecTime   time.Time
	EnergyUsed int64
	ExecResult domain.EventTaskExecResult
}
type TaskExecuteUnit struct {
	Id 			bson.ObjectId
	Task 		domain.Task				//对应Task
	DeviceId	string					//执行任务的设备
	Cpu 		int64					//分配给它的CPU
	Memory 		int64					//分配给他的内存
	Disk 		int64					//分配给他的磁盘
	NetRate 	int64					//分配给他的带宽
	State		executeState			//当前状态
	Result		result					//执行结果
}
