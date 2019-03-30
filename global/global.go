package global

import "../models"

//全局异常
var (
	Errors            []error
	SingletonNodeInfo *models.Data
	IsMaster          bool
	MasterUrl         string
	NodeStatus        int //1.正常 2.不可用
)
