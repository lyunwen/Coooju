package data

import (
	"../global"
	"../models"
	"regexp"
)

func Load() {
	dataInit()
}

func dataInit() {
	localDataInit()
	globalDataInit()
}

func globalDataInit() {
	global.SelfFlag = 1
	var err error
	global.SingletonNodeInfo = new(models.Data).GetData()
	if err != nil {
		panic(err)
	}
	global.SelfFlag = 1
	global.CuCluaster = nil
}

func localDataInit() {
	localData := new(models.Data).GetData()
	isMatch, err := regexp.Match(`^[A-Za-z0-9]{3,100}\d{1,5}$`, []byte(localData.Version))
	if err != nil {
		panic(err)
	}
	if isMatch == true {
		return
	} else {
		///备份当前的
	}
}
