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
	case 1: //初始态
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
								log.Default("初始态->备机状态")
								global.SelfFlag = 2
								global.MasterUrl = item.Address
								break
							}
						}
					}
				}
			}
			if i+1 == len(global.SingletonNodeInfo.Clusters) {
				log.Default("初始态->主机状态")
				global.SelfFlag = 3
				global.MasterUrl = global.LocalUrl
			}
		}
	case 2:
		client := &http.Client{}
		for i, item := range global.SingletonNodeInfo.Clusters {
			request, err := http.NewRequest("GET", "http://"+item.Address+"/api/IsMaster/", nil)
			if err == nil && item.Address != global.LocalUrl {
				response, err := client.Do(request)
				if err == nil && response.StatusCode == 200 {
					body, err := ioutil.ReadAll(response.Body)
					if err == nil {
						bodyStr := string(body)
						var backJsonObj ClusterBackObj
						if err = json.Unmarshal([]byte(bodyStr), &backJsonObj); err == nil {
							if backJsonObj.Code == "3" { //遇到主机则正常
								break
							}
						}
					}
				}
			}
			if i+1 == len(global.SingletonNodeInfo.Clusters) {
				log.Default("备机状态->主机状态")
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

//获取数据更新方法
func SynchronyData() error {
	if global.SelfFlag == 2 {
		client := new(http.Client)
		request, err := http.NewRequest("GET", "http://"+global.MasterUrl+"/api/getData", nil)
		if err != nil {
			return err
		}
		response, err := client.Do(request)
		if err != nil {
			return err
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		bodyStr := string(body)
		var masterData *models.Data
		if err = json.Unmarshal([]byte(bodyStr), &masterData); err != nil {
			return err
		}
		var localData = new(models.Data).GetData()
		localPre, localVersion := localData.GetVersionInfo()
		masterPre, masterVersion := masterData.GetVersionInfo()
		if localPre == masterPre {
			if localVersion <= masterVersion {
				masterData.SetData()
				return nil
			} else {
				return errors.New("local version bigger than master version")
			}
		} else {
			return errors.New("local pre not equal master pre")
		}
	} else {
		return errors.New("only SelfFlag=2 can be SynchronyData")
	}
}
