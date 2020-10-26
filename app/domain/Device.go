package domain

import (
	"gopkg.in/mgo.v2/bson"
)

type Device struct {
	Id          	bson.ObjectId			`bson:"_id,omitempty" json:"id"`				//Device本身在数据库里的ID
	EdgexId			string					`json:"edgexid"`								//Device在Edgex里的ID
	Name			string					`json:"name"`									//device的名称
	Cpu				int64					`json:"cpu"`									//device的CPU
	Memory			int64					`json:"memory"`									//device的内存容量
	Disk			int64					`json:"disk"`									////device的硬盘容量
	GetCommands		map[string]*Command		`json:"getcommands"`							//device的command,get
	PutCommands		map[string]*Command		`json:"putcommands"`							//device的command,put
}
