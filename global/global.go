package global

import "../models"

//全局异常
var (
	Errors            []error
	SingletonNodeInfo *models.Data
	//MasterFlag        int // -1:异常态1:初始态 2：正常态 3：异常态
	MasterUrl  string
	SelfFlag   int // -1:异常态 1：初始态 2：备机状态 3：主机状态
	NodeStatus int //1.正常 2.不可用
)
