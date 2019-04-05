package data

import (
	"../global"
	"../models"
)

func Load() error {
	var err error
	err = dataInit()
	return err
}

func dataInit() error {
	var err error
	global.MasterFlag = 1
	global.SingletonNodeInfo, err = new(models.Data).GetData()
	if err != nil {
		return err
	}
	global.NodeStatus = 2
	return err
}
