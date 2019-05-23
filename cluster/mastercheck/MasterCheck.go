package mastercheck

import (
	"../../cluster"
	"../../common/log"
	"../../global"
	"../../global/SelfFlagStatus"
	"../../models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func Check() error {
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	// 类raft算法
	switch global.SelfFlag {
	case SelfFlagStatus.Init:
		for i, item := range global.SingletonNodeInfo.Clusters {
			code, info, err := checkUrl(client, item.Address)
			if err != nil {
				log.Warn("checkUrl error:" + err.Error())
			} else if code == strconv.Itoa(int(SelfFlagStatus.Leader)) {
				if err := cluster.SynchronyData(info.Address); err != nil {
					log.Warn("初始态->异常态 数据同步不成功:" + err.Error())
					panic("初始态->异常态 数据同步不成功:" + err.Error())
				} else {
					log.Warn("初始态->备机状态  Master Address:" + info.Address)
					global.SelfFlag = SelfFlagStatus.Follow
					global.MasterUrl = info.Address
				}
				break
			}
			if i+1 == len(global.SingletonNodeInfo.Clusters) {
				global.MasterUrl = global.CuCluster.Address
				global.SelfFlag = SelfFlagStatus.Follow
				log.Warn("初始态->主机状态")
			}
		}
	case SelfFlagStatus.Follow:
		code, _, err := checkUrl(client, global.MasterUrl)
		if err != nil {
			log.Warn("check master Url error:" + err.Error())
		}
		if code == strconv.Itoa(int(SelfFlagStatus.Leader)) {
			return nil
		} else {
			for i, item := range global.SingletonNodeInfo.Clusters {
				code, info, err := checkUrl(client, item.Address)
				if err != nil {
					log.Warn("check Url error:" + err.Error())
				} else if code == strconv.Itoa(int(SelfFlagStatus.Leader)) {
					if global.MasterUrl != info.Address {
						log.Warn("主机已切换: " + global.MasterUrl + "->" + info.Address + "")
						global.MasterUrl = info.Address
					}
					return nil
				}
				if i+1 == len(global.SingletonNodeInfo.Clusters) {
					log.Warn("备机状态->主机状态")
					global.SelfFlag = SelfFlagStatus.Leader
					global.MasterUrl = global.CuCluster.Address
				}
			}
		}
	case SelfFlagStatus.Candidate:
		hasVotedCount := 0
		for _, item := range global.SingletonNodeInfo.Clusters {
			fmt.Print(item)
		}
		if hasVotedCount > len(global.SingletonNodeInfo.Clusters) {

		}
	case SelfFlagStatus.Leader:
		for _, item := range global.SingletonNodeInfo.Clusters {
			code, info, err := checkUrl(client, item.Address)
			if err != nil {
				log.Warn("check Url error:" + err.Error())
			} else if code == strconv.Itoa(int(SelfFlagStatus.Leader)) {
				if info.Level > global.CuCluster.Level {
					log.Warn("当前master【address:" + global.CuCluster.Address + " level:" + strconv.Itoa(global.CuCluster.Level) + " name:" + global.CuCluster.Name + "】 发现 另外 master【address:" + info.Address + " level:" + strconv.Itoa(info.Level) + " name:" + info.Name + "】")
					global.SelfFlag = 2
					global.MasterUrl = info.Address
					log.Warn("主机状态->备机状态 Master Address:" + info.Address + " (遇到权重更好主机)")
				}
			}
		}
	default:
		panic("当前机器状态" + strconv.Itoa(int(global.SelfFlag)) + "异常 停止检测")
	}
	return nil
}

func checkUrl(client *http.Client, url string) (string, *models.Cluster, error) {
	request, err := http.NewRequest("GET", "http://"+url+"/api/IsMaster/", nil)
	if err != nil {
		return "", nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return "", nil, err
	}
	if response.StatusCode != 200 {
		return "", nil, errors.New("StatusCode error")
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", nil, err
	}
	bodyStr := string(body)
	var dataMsg json.RawMessage
	var backJsonObj = cluster.BackObj{Data: &dataMsg}
	err = json.Unmarshal([]byte(bodyStr), &backJsonObj)
	if err != nil {
		return "", nil, err
	}
	var otherMaster *models.Cluster
	err = json.Unmarshal(dataMsg, &otherMaster)
	if err != nil {
		return "", nil, err
	}
	return backJsonObj.Code, otherMaster, err
}
