package cluster

import (
	"../global"
	"../models"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

type BackObj struct {
	Code string
	Msg  string
	Data interface{}
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
				global.CuCluster = &item
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

//获取master数据更新本地
func SynchronyData(url string) error {
	client := new(http.Client)
	request, err := http.NewRequest("GET", "http://"+url+"/api/cluster/getData", nil)
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
	var msg json.RawMessage
	var returnObj = &BackObj{
		Data: &msg,
	}
	if err := json.Unmarshal([]byte(bodyStr), &returnObj); err != nil {
		return err
	}
	var masterData *models.Data
	if err = json.Unmarshal(msg, &masterData); err != nil {
		return err
	}
	err = masterData.SetData()
	return err
}
