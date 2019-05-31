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
		otherNode, err := getNodeInfo(client, global.CurrentData.MasterAddress)
		if err != nil {
			log.Warn("[Cluster State]: Follow->Candidate 连接主机URL:" + global.CurrentData.Address + "异常")
			global.CurrentData.ClusterState = clusterState.Candidate
			break
		}
		if otherNode.ClusterState != clusterState.Leader {
			log.Warn("[Cluster State]: Follow->Candidate 获取主机URL:" + global.CurrentData.Address + "状态" + strconv.Itoa(int(otherNode.ClusterState)) + "异常")
			global.CurrentData.ClusterState = clusterState.Candidate
			break
		}
		dataInfo, err := getCusterInfo(client, global.CurrentData.Address)
		if err != nil {
			log.Warn("[Cluster State]: Follow->Candidate 获取集群数据URL:" + global.CurrentData.Address + "异常")
			global.CurrentData.ClusterState = clusterState.Candidate
			break
		}
		err = dataInfo.SetData()
		if err != nil {
			log.Warn("[Cluster State]: Follow->Candidate SetData Error" + err.Error())
			global.CurrentData.ClusterState = clusterState.Candidate
			break
		}
	case clusterState.Candidate: //1.拉票成主机
		var moreThanHalf = len(global.ClusterData.Services) / 2
		var connectCount = 0
		var votedCount = 0
		for _, item := range global.ClusterData.Clusters {
			votedResult, err := getVotes(client, item.Address, strconv.Itoa(global.CurrentData.VotedTerm))
			if err == nil {
				if votedResult == "found leader" {
					log.Warn("[Cluster State]: Candidate->Follow found leader:" + item.Address)
					global.CurrentData.ClusterState = clusterState.Follow
					return nil
				} else if votedResult == "ok" {
					votedCount++
					log.Warn("拉票成功 URL:" + item.Address + " Term:" + strconv.Itoa(global.CurrentData.VotedTerm) + " Votes:" + strconv.Itoa(votedCount))
				}
				connectCount++
			} else {
				log.Warn("拉票失败 URL:" + item.Address + " Term:" + strconv.Itoa(global.CurrentData.VotedTerm) + " Error:" + err.Error())
			}
		}
		if connectCount > moreThanHalf {
			global.CurrentData.VotedTerm++
		}
		if votedCount > moreThanHalf {
			log.Warn("[Cluster State]: Candidate->Leader Votes:" + strconv.Itoa(votedCount))
			global.CurrentData.ClusterState = clusterState.Leader
			global.CurrentData.Term = global.CurrentData.VotedTerm
			break
		}
	case clusterState.Leader: // 1.检查是否存在更高优先级主机
		for _, item := range global.ClusterData.Clusters {
			otherNode, err := getNodeInfo(client, item.Address)
			if err != nil {
				log.Warn("连接主机URL:" + item.Address + "异常")
				break
			}
			if otherNode.ClusterState == clusterState.Leader && global.CurrentData.VotedTerm < otherNode.VotedTerm {
				log.Warn("Found Bigger Term: Local（URL:" + global.CurrentData.Address + " Term：" + strconv.Itoa(global.CurrentData.VotedTerm) + "） Other（URL:" + otherNode.Address + " Term:" + strconv.Itoa(otherNode.VotedTerm) + "）")
				global.CurrentData.VotedTerm = otherNode.VotedTerm
				global.CurrentData.ClusterState = clusterState.Follow
				log.Warn("[Cluster State]: Leader->Follow")
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

func getVotes(client *http.Client, url string, term string) (string, error) {
	request, err := http.NewRequest("GET", "http://"+url+"/api/setVotes?term="+term, nil)
	if err != nil {
		return "", err
	}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	if response.StatusCode != 200 {
		return "", errors.New("StatusCode error")
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	bodyStr := string(body)
	var dataMsg json.RawMessage
	var backJsonObj = cluster.ClusterBackObj{Data: &dataMsg}
	err = json.Unmarshal([]byte(bodyStr), &backJsonObj)
	if err != nil {
		return "", err
	}
	if backJsonObj.Code == "0" {
		var resultData string
		if err = json.Unmarshal(dataMsg, &resultData); err != nil {
			return "", err
		}
		return resultData, nil
	} else {
		return "", errors.New(backJsonObj.Msg)
	}
}
