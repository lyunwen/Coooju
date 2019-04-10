package timer

import (
	"../cluster"
	"../common/log"
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
			err := cluster.SynchronyData()
			if err != nil {
				log.Error(err.Error())
			}
		}
		ch <- 1
	}()
	<-ch
}
