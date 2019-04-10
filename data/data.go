package data

import (
	"../common/log"
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

//运行程序数据初始化
func globalDataInit() {
	global.SelfFlag = 1
	var err error
	global.SingletonNodeInfo = new(models.Data).GetData()
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
	global.SelfFlag = 1
	global.CuCluaster = nil
}

//data.json 格式校验
func localDataInit() {
	localData := new(models.Data).GetData()
	if localData.Version == "" {
		return
	}
	isMatch, err := regexp.Match(`^[A-Za-z0-9]{3,100}-[\d]{1,5}$`, []byte(localData.Version))
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
	if isMatch == true {
		return
	} else {
		log.Error("version format error")
		panic(" version:" + localData.Version + " format error")
	}
}
