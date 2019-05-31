package mastercheck

import (
	"../../cluster"
	"../../common/log"
	"../../global"
	"../../models"
	"../../models/clusterState"
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
	switch global.CurrentData.ClusterState {
	case clusterState.Follow: //1.同步主机数据 2.转候选者
		for i, item := range global.ClusterData.Clusters {
			if item.State == clusterState.Leader {
				otherNode, err := getNodeInfo(client, item.Address)
				if err != nil {
					log.Warn("连接主机URL:" + item.Address + "异常")
					break
				}
				if otherNode.ClusterState != clusterState.Leader {
					log.Warn("获取主机URL:" + item.Address + "状态" + strconv.Itoa(int(otherNode.ClusterState)) + "异常")
					break
				}
				dataInfo, err := getCusterInfo(client, item.Address)
				if err != nil {
					log.Warn("连接主机URL:" + item.Address + "异常")
					break
				}
				err = dataInfo.SetData()
				if err != nil {
					log.Warn("SetData Error" + err.Error())
					break
				}
			}
			if i+1 == len(global.ClusterData.Clusters) {
				global.CurrentData.ClusterState = clusterState.Candidate
				log.Warn("State Follow->Candidate 未发现主机")
			}
		}
	case clusterState.Candidate: //1.拉票成主机
		var votedCount = 0
		global.CurrentData.VotedTerm++
		for _, item := range global.ClusterData.Clusters {
			err := getVotes(client, item.Address, strconv.Itoa(global.CurrentData.VotedTerm))
			if err == nil {
				log.Warn("拉票成功 URL:" + item.Address + " Term:" + strconv.Itoa(global.CurrentData.VotedTerm))
				votedCount++
			} else {
				log.Warn("拉票失败 URL:" + item.Address + " Term:" + strconv.Itoa(global.CurrentData.VotedTerm) + " Error:" + err.Error())
			}
			if votedCount > (len(global.ClusterData.Services) / 2) {
				log.Warn("State Candidate->Leader 获取足够票数")
				global.CurrentData.ClusterState = clusterState.Leader
				break
			}
		}
	case clusterState.Leader: // 1.检查是否存在更高优先级主机
		for _, item := range global.ClusterData.Clusters {
			otherNode, err := getNodeInfo(client, item.Address)
			if err != nil {
				log.Warn("连接主机URL:" + item.Address + "异常")
				break
			}
			if otherNode.ClusterState == clusterState.Leader && global.CurrentData.VotedTerm < otherNode.VotedTerm {
				log.Warn("发现权重更高Leader: 本机Leader（URL:" + global.CurrentData.Address + " Term：" + strconv.Itoa(global.CurrentData.VotedTerm) + "） 检测Leader（URL:" + otherNode.Address + " Term:" + strconv.Itoa(otherNode.VotedTerm) + "）")
				global.CurrentData.VotedTerm = otherNode.VotedTerm
				global.CurrentData.ClusterState = clusterState.Follow
				log.Warn("Leader->Follow")
				break
			}
		}
	default:
		log.Error("当前机器状态" + strconv.Itoa(int(global.CurrentData.ClusterState)) + "异常 停止检测")
		return nil
	}
	return nil
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
	var backJsonObj = cluster.ClusterBackObj{Data: &dataMsg}
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

func getNodeInfo(client *http.Client, url string) (*global.CurrentNodeInfo, error) {
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
	var backJsonObj = cluster.ClusterBackObj{Data: &dataMsg}
	err = json.Unmarshal([]byte(bodyStr), &backJsonObj)
	if err != nil {
		return nil, err
	}
	if backJsonObj.Code != "0" {
		return nil, errors.New("get data error")
	}
	var otherNode *global.CurrentNodeInfo
	err = json.Unmarshal(dataMsg, &otherNode)
	if err != nil {
		return nil, err
	}
	return otherNode, err
}

func getVotes(client *http.Client, url string, term string) error {
	request, err := http.NewRequest("GET", "http://"+url+"/api/setVotes?term="+term, nil)
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New("StatusCode error")
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	bodyStr := string(body)
	var dataMsg json.RawMessage
	var backJsonObj = cluster.ClusterBackObj{Data: &dataMsg}
	err = json.Unmarshal([]byte(bodyStr), &backJsonObj)
	if err != nil {
		return err
	}
	//var resultData string
	//if err = json.Unmarshal(dataMsg, &resultData); err != nil {
	//	return "", err
	//}
	if backJsonObj.Code == "0" {
		return nil
	} else {
		return errors.New(backJsonObj.Msg)
	}
}
