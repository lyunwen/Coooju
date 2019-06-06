package cluster

import (
	"../cluster/clusterState"
	"../models"
	"net"
	"strings"
)

var (
	ClusterData *models.Data
	CurrentData *CurrentNodeInfo
)

type CurrentNodeInfo struct {
	MasterAddress string
	ClusterState  clusterState.ClusterState
	Name          string
	Level         int
	Address       string
}

func Init() {
	ClusterData = new(models.Data).GetData()
	CurrentData = new(CurrentNodeInfo)
	for index, item := range ClusterData.Clusters {
		isLocalIp, err := isLocalIp(strings.Split(item.Address, ":")[0])
		if err != nil {
			panic(err.Error())
		}
		if isLocalIp {
			conn, err := net.Dial("tcp", item.Address)
			if err != nil {
				CurrentData.ClusterState = clusterState.Follow
				CurrentData.Address = item.Address
				CurrentData.Level = item.Level
				CurrentData.MasterAddress = ""
			}
			_ = conn.Close()
		}
		if index == len(ClusterData.Clusters) {
			panic("can not find right ip")
		}
	}
}

func isLocalIp(ip string) (bool, error) {
	adders, err := net.InterfaceAddrs()
	if err != nil {
		return false, err
	}
	for _, address := range adders {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
				return true, nil
			}
		}
	}
	return false, nil
}
