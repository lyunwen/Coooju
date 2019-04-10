package timer

import (
	"../cluster"
	"time"
)

func Load() {
	MasterCheckLoop()
	SynchronyDataLoop()
}

//主机检测
func MasterCheckLoop() {
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

func SynchronyDataLoop() {
	var ch chan int
	//定时任务
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for range ticker.C {
			cluster.SynchronyData()
		}
		ch <- 1
	}()
	<-ch
}
