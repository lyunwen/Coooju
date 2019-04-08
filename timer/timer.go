package timer

import (
	"../cluster"
	"time"
)

func Load() {
	MasterCheck()
	GetClusterData()
}

//主机检测
func MasterCheck() {
	var ch chan int
	//定时任务
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for range ticker.C {
			cluster.MasterCheck()
		}
		ch <- 1
	}()
	<-ch
}

func GetClusterData() {
	var ch chan int
	//定时任务
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for range ticker.C {
			cluster.GetClusterData()
		}
		ch <- 1
	}()
	<-ch
}
