package timer

import (
	"../cluster"
	"time"
)

func Load() {
	MasterCheck()
	SynchronyNodeData()
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

}

//数据同步
func SynchronyNodeData() {
	var ch chan int
	//定时任务
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for range ticker.C {
			cluster.SynchronyNodeData()
		}
		ch <- 1
	}()
	<-ch
}
