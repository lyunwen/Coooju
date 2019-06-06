package mastercheck

import (
	"../../cluster"
	"../../common/log"
	"../../models"
	"../clusterState"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func Check() error {
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	switch cluster.OwnData.ClusterState {
	case clusterState.Follow:
		masterNode, err := getNodeInfo(client, cluster.OwnData.MasterAddress)
		if err == nil && masterNode.ClusterState == clusterState.Leader {
			dataInfo, err := getCusterInfo(client, cluster.OwnData.Address)
			if err != nil {
				log.Warn("[Follow]get master node URL:" + cluster.OwnData.Address + "error")
				break
			}
			if err = dataInfo.SetData(); err != nil {
				log.Warn("[Follow]set data error" + err.Error())
				break
			}
			for _, item := range dataInfo.Clusters {
				if item.Address == cluster.OwnData.Address {
				}
			}
		} else {
			for _, item := range cluster.ShareData.Clusters {
				otherNode, err := getNodeInfo(client, item.Address)
				if err == nil && otherNode.ClusterState == clusterState.Leader {
					log.Warn("[Follow]master exchange:" + cluster.OwnData.MasterAddress + " ->" + item.Address)
					cluster.OwnData.MasterAddress = item.Address
					return nil
				}
			}
			log.Warn("[Follow]become master:" + cluster.OwnData.MasterAddress + " ->" + cluster.OwnData.Address)
			cluster.OwnData.ClusterState = clusterState.Leader
			cluster.OwnData.MasterAddress = cluster.OwnData.Address
		}
	case clusterState.Leader:
		for _, item := range cluster.ShareData.Clusters {
			node, err := getNodeInfo(client, item.Address)
			if err != nil {
				log.Warn("[Leader]Url:" + item.Address + "check error：" + err.Error())
			} else if node.ClusterState == clusterState.Follow {
				log.Warn("[Leader]Url:" + item.Address + "check success")
			} else if node.ClusterState == clusterState.Leader && node.GetLevel() > cluster.OwnData.GetLevel() {
				log.Warn("[Leader]Url:" + item.Address + "find bigger level:" + strconv.Itoa(node.GetLevel()) + ">" + strconv.Itoa(cluster.OwnData.GetLevel()))
				log.Warn("[Leader]master exchange:" + cluster.OwnData.MasterAddress + " ->" + node.Address)
				cluster.OwnData.MasterAddress = node.Address
			} else {
				log.Warn("Url:" + item.Address + " State:" + strconv.Itoa(int(clusterState.Leader)))
			}
		}
	default:
		panic("当前机器状态" + strconv.Itoa(int(cluster.OwnData.ClusterState)) + "异常 停止检测")
	}
	return nil
}

func getNodeInfo(client *http.Client, url string) (*cluster.CurrentNodeInfo, error) {
	request, err := http.NewRequest("GET", "http://"+url+"/api/cluster/getNodeInfo/", nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New("StatusCode error")
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	bodyStr := string(body)
	var dataMsg json.RawMessage
	var backJsonObj = cluster.BackObj{Data: &dataMsg}
	err = json.Unmarshal([]byte(bodyStr), &backJsonObj)
	if err != nil {
		return nil, err
	}
	if backJsonObj.Code != "0" {
		return nil, errors.New("get data error")
	}
	var otherNode *cluster.CurrentNodeInfo
	err = json.Unmarshal(dataMsg, &otherNode)
	if err != nil {
		return nil, err
	}
	return otherNode, err
}

func getCusterInfo(client *http.Client, url string) (*models.Data, error) {
	request, err := http.NewRequest("GET", "http://"+url+"/api/cluster/getCusterInfo/", nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New("StatusCode error")
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	bodyStr := string(body)
	var dataMsg json.RawMessage
	var backJsonObj = cluster.BackObj{Data: &dataMsg}
	err = json.Unmarshal([]byte(bodyStr), &backJsonObj)
	if err != nil {
		return nil, err
	}
	if backJsonObj.Code != "0" {
		return nil, errors.New("get data error")
	}
	var otherMaster *models.Data
	err = json.Unmarshal(dataMsg, &otherMaster)
	if err != nil {
		return nil, err
	}
	return otherMaster, err
}
