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
	switch cluster.CurrentData.ClusterState {
	case clusterState.Follow:
		masterNode, err := getNodeInfo(client, cluster.CurrentData.MasterAddress)
		if err == nil && masterNode.ClusterState == clusterState.Leader {
			dataInfo, err := getCusterInfo(client, cluster.CurrentData.Address)
			if err != nil {
				log.Warn("[Follow]get master node URL:" + cluster.CurrentData.Address + "error")
				break
			}
			if err = dataInfo.SetData(); err != nil {
				log.Warn("[Follow]set data error" + err.Error())
				break
			}
		} else {
			for _, item := range cluster.ClusterData.Clusters {
				otherNode, err := getNodeInfo(client, item.Address)
				if err == nil && otherNode.ClusterState == clusterState.Leader {
					log.Warn("[Follow]master exchange:" + cluster.CurrentData.MasterAddress + " ->" + item.Address)
					cluster.CurrentData.MasterAddress = item.Address
					break
				}
			}
			cluster.CurrentData.ClusterState = clusterState.Leader
			cluster.CurrentData.MasterAddress = cluster.CurrentData.Address
			log.Warn("[Follow]become master:" + cluster.CurrentData.MasterAddress + " ->" + cluster.CurrentData.Address)
		}
	case clusterState.Leader:
		for _, item := range cluster.ClusterData.Clusters {
			node, err := getNodeInfo(client, item.Address)
			if err != nil {
				log.Warn("[Leader]Url:" + item.Address + "check error：" + err.Error())
			} else if node.ClusterState == clusterState.Follow {
				log.Warn("[Leader]Url:" + item.Address + "check success")
			} else if node.ClusterState == clusterState.Leader && node.Level > cluster.CurrentData.Level {
				log.Warn("[Leader]Url:" + item.Address + "find bigger level:" + strconv.Itoa(node.Level) + ">" + strconv.Itoa(cluster.CurrentData.Level))
				log.Warn("[Leader]master exchange:" + cluster.CurrentData.MasterAddress + " ->" + node.Address)
				cluster.CurrentData.MasterAddress = node.Address
			} else {
				log.Warn("Url:" + item.Address + "error")
			}
		}
	default:
		panic("当前机器状态" + strconv.Itoa(int(cluster.CurrentData.ClusterState)) + "异常 停止检测")
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
