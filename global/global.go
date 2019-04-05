package global

import "../models"

//全局异常
var (
	Errors            []error
	SingletonNodeInfo *models.Data
	MasterFlag        int // -1:非正常状态，手动处理 1：待转备机或升主机 2：备机状态 3：主机状态
	MasterUrl         string
	NodeStatus        int //1.正常 2.不可用
)
