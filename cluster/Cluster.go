package cluster

import (
	"../common"
	"../global"
	"../models"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func MasterCheck() {
	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", "http://127.0.0.1/health", nil)

	reqest.Header.Set("Accept", "application/json")
	reqest.Header.Set("Accept-Charset", "utf-8")
	reqest.Header.Set("Connection", "keep-alive")

	response, err := client.Do(reqest)
	if err != nil {
		global.IsMaster = true
	}
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		bodystr := string(body)
		fmt.Println(bodystr)
	}
}

func SynchronyNodeData() {
	for _, item := range global.SingletonNodeInfo.Clusters {
		result, err := synchronyNodeData(global.SingletonNodeInfo, item.Address+"/api/SynchronyNodeData")
		if err != nil {
			_ = common.Log(time.Now().Format("2005-01-02 15:04:05") + "synchronyNodeData error:" + err.Error())
		}
		switch result {
		case "ok":
			return
		case "smaller":
			_ = common.Log(time.Now().Format("2005-01-02 15:04:05") + "smaller smaller error:" + err.Error())
			return
		case "equal":
			_ = common.Log(time.Now().Format("2005-01-02 15:04:05") + "equal version error:" + err.Error())
			return
		default:
			_ = common.Log(time.Now().Format("2005-01-02 15:04:05") + "synchronyNodeData result error:" + err.Error())
			return
		}
	}
}

func synchronyNodeData(data *models.Data, url string) (string, error) {
	client := new(http.Client)
	dataJsonByte, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataJsonStr := string(dataJsonByte)
	request, err := http.NewRequest("GET", url+"?data="+dataJsonStr, nil)
	response, err := client.Do(request)
	if err != nil {
		global.IsMaster = true
	}
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		bodystr := string(body)
		fmt.Println(bodystr)
	}
	return "", nil
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
