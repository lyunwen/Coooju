package global

import (
	"../models"
	"./SelfFlagStatus"
)

//全局异常
var (
	Errors            []error
	SingletonNodeInfo *models.Data
	MasterUrl         string
	CuCluster         *models.Cluster
	SelfFlag          SelfFlagStatus.SelfFlagStatus
)
