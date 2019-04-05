package cluster

import (
	"../common"
	"../global"
	"../models"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"net/http"
	url2 "net/url"
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
		common.Log("服务状态调整 当前状态：初始态" + strconv.Itoa(global.SelfFlag))
	case 1:
		fallthrough //初始态
	case 2:
		common.Log("服务状态调整 当前状态：初始态" + strconv.Itoa(global.SelfFlag))
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
								common.Log("切备机")
								global.SelfFlag = 2
								global.MasterUrl = item.Address
								break
							}
						}
					}
				}
			}
			if i+1 == len(global.SingletonNodeInfo.Clusters) {
				common.Log("切主机")
				global.SelfFlag = 3
				global.MasterUrl = global.LocalUrl
			}
		}
	case 3: //主机状态
	default:
		global.SelfFlag = -1
		common.Log("当前机器状态" + strconv.Itoa(global.SelfFlag) + "异常 停止检测")
	}

}

func SynchronyNodeData() {
	for _, item := range global.SingletonNodeInfo.Clusters {
		synchronyNodeData(global.SingletonNodeInfo, item.Address)
		//if err != nil {
		//	_ = common.Log(time.Now().Format("2005-01-02 15:04:05") + "synchronyNodeData error:" + err.Error())
		//}
		//switch result {
		//case "ok":
		//	return
		//case "smaller":
		//	_ = common.Log(time.Now().Format("2005-01-02 15:04:05") + "smaller smaller error:" + err.Error())
		//	return
		//case "equal":
		//	_ = common.Log(time.Now().Format("2005-01-02 15:04:05") + "equal version error:" + err.Error())
		//	return
		//default:
		//	_ = common.Log(time.Now().Format("2005-01-02 15:04:05") + "synchronyNodeData result error:" + err.Error())
		//	return
		//}
	}
}

// 地址/状态 1 可用 2网络不可用 4系统不可用 5系统内部未知异常 6系统内部已知异常 7同步成功 8无需同步 9不能同步
var State = map[string]int{}

func synchronyNodeData(data *models.Data, url string) {
	client := new(http.Client)
	dataJsonBytes := make([]byte, 1024)
	dataJsonBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	dataJsonStr := string(dataJsonBytes)
	request, err := http.NewRequest("GET", "http://"+url+"/api/SynchronyNodeData?data="+url2.PathEscape(dataJsonStr), nil)
	if err != nil {
		panic(err)
	}
	response, err := client.Do(request)
	if err != nil {
		State[url] = 2
		return
	}
	if response.StatusCode != 200 {
		State[url] = 5
		return
	}
	body, _ := ioutil.ReadAll(response.Body)
	bodyStr := string(body)
	var backObj ClusterBackObj
	if err := json.Unmarshal([]byte(bodyStr), &backObj); err != nil {
		State[url] = 5
		return
	}
	if backObj.Code != "0" {
		State[url] = 6
		return
	}
	State[url] = 7
	return
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
		request, err := http.NewRequest("GET", "http://"+global.MasterUrl+"/api/SynchronyNodeData", nil)
		if err != nil {
			panic(err)
		}
		response, err := client.Do(request)
		if err != nil {
			return
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			///
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
