package timer

import (
	"../cluster"
	"../common/log"
	"../global"
	"time"
)

func Load() {
	go MasterCheckLoop()
	go SynchronyDataLoop()
}

//主机检测
func MasterCheckLoop() {
	var ch chan int
	//定时任务
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for range ticker.C {
			if err := cluster.MasterCheck(); err != nil {
				log.Error(err.Error())
			}
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
			if global.SelfFlag == 2 {
				if err := cluster.SynchronyData(); err != nil {
					log.Error(err.Error())
				}
			}
		}
		ch <- 1
	}()
	<-ch
}
