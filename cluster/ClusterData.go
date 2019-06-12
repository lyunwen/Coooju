package cluster

import (
	"../cluster/clusterState"
	"../models"
	"net"
	"strings"
)

var (
	ShareData *models.Data
	OwnData   *CurrentNodeInfo
)

type CurrentNodeInfo struct {
	MasterAddress string
	ClusterState  clusterState.ClusterState
	Address       string
	Name          string
	Level         int
}

func (info *CurrentNodeInfo) GetName() string {
	for _, item := range ShareData.Clusters {
		if item.Address == info.Address {
			return item.Name
		}
	}
	panic("no find")
}

func (info *CurrentNodeInfo) GetLevel() int {
	for _, item := range ShareData.Clusters {
		if item.Address == info.Address {
			return item.Level
		}
	}
	panic("no find")
}

func Init() {
	ShareData = new(models.Data).GetData()
	for index, item := range ShareData.Clusters {
		if isLocalIP(strings.Split(item.Address, ":")[0]) {
			conn, err := net.Dial("tcp", item.Address)
			if err == nil {
				_ = conn.Close()
			} else {
				OwnData = &CurrentNodeInfo{
					ClusterState:  clusterState.Follow,
					Address:       item.Address,
					MasterAddress: "",
				}
				break
			}
		}
		if index+1 == len(ShareData.Clusters) {
			panic("can not find right ip")
		}
	}
}

func isLocalIP(ip string) bool {
	adders, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, value := range adders {
		if ipNet, ok := value.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				if ipNet.IP.String() == ip {
					return true
				}
			}
		}
	}
	return false
}
