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
	code string
	msg  string
	data interface{}
}

//非线程安全
func MasterCheck() {
	//-1:非正常状态，手动处理 1：待转备机或升主机 2：备机状态 3：主机状态
	switch global.MasterFlag {
	case -1: //异常状态不检查
		common.Log("当前机器状态" + strconv.Itoa(global.MasterFlag) + "异常 停止检测")
		break
	case 1: //修整
		common.Log("当前机器状态待修整")

		break
	case 2:
		client := &http.Client{}
		for _, item := range global.SingletonNodeInfo.Clusters {
			request, err := http.NewRequest("GET", "http://"+item.Address+"/api/IsMaster/", nil)
			if err != nil {
				common.Log("获取 地址：" + item.Address + "error1")
				continue
			} else {
				response, err := client.Do(request)
				if err != nil {
					common.Log("获取 地址：" + item.Address + "error2")
					continue
				}
				if response.StatusCode != 200 {
					common.Log("获取 地址：" + item.Address + "error3")
					continue
				} else {
					body, err := ioutil.ReadAll(response.Body)
					if err != nil {
						common.Log("获取 地址：" + item.Address + "error4")
						continue
					} else {
						backJsonStr := string(body)
						var backJsonObj ClusterBackObj
						if json.Unmarshal([]byte(backJsonStr), &backJsonObj) != nil {
							var dataStr = backJsonObj.data.(map[string]interface{})["description"].(string)
							if dataStr == "3" { //遇到主机切备机
								global.MasterFlag = 2
								break
							} else if dataStr == "2" { //遇到备机往下走
								global.MasterFlag = -1
							} else { //其他情况切异常
								global.MasterFlag = -1
								common.Log("获取 地址：" + item.Address + "节点信息失败5")
								break
							}
						}
					}
				}
			}
		}
	case 3: //自己是主机停止检查
		break
	default:
		global.MasterFlag = -1
		common.Log("当前机器状态" + strconv.Itoa(global.MasterFlag) + "异常 停止检测")
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
	if backObj.code != "0" {
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
