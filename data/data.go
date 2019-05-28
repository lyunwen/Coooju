package data

import (
	"../global"
	"../models"
)

func Load() {
	dataInit()
}

func dataInit() {
	global.ClusterData = new(models.Data).GetData()
}
