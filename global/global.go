package global

import "../models"

//全局异常
var (
	Errors            []error
	SingletonNodeInfo *models.Data
	MasterUrl         string
	CuCluster         *models.Cluster
	SelfFlag          int // -1:异常态 1：初始态 2：备机状态 3：主机状态
)
