package data

import (
	"../global"
	"../models"
)

func Load() {
	dataInit()
}

func dataInit() {
	global.SelfFlag = 1
	var err error
	global.SingletonNodeInfo = new(models.Data).GetData()
	if err != nil {
		panic(err)
	}
	global.SelfFlag = 1
	global.CuCluaster = nil
}
