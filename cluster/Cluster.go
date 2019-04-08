package cluster

import (
	"../common/log"
	"../global"
	"../models"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type ClusterBackObj struct {
	Code string
	Msg  string
	Data interface{}
}

//非线程安全 服务状态调整
func MasterCheck() {
	// -1:异常态 1：初始态 2：备机状态 3：主机状态
	switch global.SelfFlag {
	case -1: //异常态
		log.Default("服务状态调整 当前状态：初始态" + strconv.Itoa(global.SelfFlag))
	case 1:
		fallthrough //初始态
	case 2:
		log.Default("服务状态调整 当前状态：初始态" + strconv.Itoa(global.SelfFlag))
		client := &http.Client{}
		for i, item := range global.SingletonNodeInfo.Clusters {
			request, err := http.NewRequest("GET", "http://"+item.Address+"/api/IsMaster/", nil)
			if err == nil {
				response, err := client.Do(request)
				if err == nil && response.StatusCode == 200 {
					body, err := ioutil.ReadAll(response.Body)
					if err == nil {
						bodyStr := string(body)
						var backJsonObj ClusterBackObj
						if err = json.Unmarshal([]byte(bodyStr), &backJsonObj); err == nil {
							if backJsonObj.Code == "3" { //遇到主机切备机
								log.Default("切备机")
								global.SelfFlag = 2
								global.MasterUrl = item.Address
								break
							}
						}
					}
				}
			}
			if i+1 == len(global.SingletonNodeInfo.Clusters) {
				log.Default("切主机")
				global.SelfFlag = 3
				global.MasterUrl = global.LocalUrl
			}
		}
	case 3: //主机状态
	default:
		global.SelfFlag = -1
		log.Error("当前机器状态" + strconv.Itoa(global.SelfFlag) + "异常 停止检测")
	}

}

func GetAvailablePortAddress() (string, error) {
	for _, item := range global.SingletonNodeInfo.Clusters {
		isLocalIp, err := isLocalIp(strings.Split(item.Address, ":")[0])
		if err != nil {
			return "", err
		}
		if isLocalIp {
			conn, err := net.Dial("tcp", item.Address)
			if err != nil {
				return item.Address, nil
			}
			_ = conn.Close()
		}
	}
	return "", errors.New("not find")
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

func GetClusterData() {
	if global.SelfFlag == 2 { //备机才能同步
		client := new(http.Client)
		request, err := http.NewRequest("GET", "http://"+global.MasterUrl+"/api/getData", nil)
		if err != nil {
			log.Error("拉取数据异常:" + err.Error())
			return
		}
		response, err := client.Do(request)
		if err != nil {
			log.Error("拉取数据异常:" + err.Error())
			return
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Error("拉取数据异常:" + err.Error())
			return
		}
		bodyStr := string(body)
		var dataObj *models.Data
		if err = json.Unmarshal([]byte(bodyStr), &dataObj); err != nil {
			global.SelfFlag = -1
			return
		} else if _, err := dataObj.SetData(); err != nil {
			global.SelfFlag = -1
			return
		}
	}
}
