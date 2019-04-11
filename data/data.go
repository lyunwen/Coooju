package data

import (
	"../global"
	"../models"
)

func Load() {
	dataInit()
}

func dataInit() {
	global.SingletonNodeInfo = new(models.Data).GetData()
	global.SelfFlag = 1
	global.CuCluster = nil
}
