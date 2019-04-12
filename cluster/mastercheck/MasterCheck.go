package mastercheck

import (
	"../../cluster"
	"../../common/log"
	"../../global"
	"../../models"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

func Check() error {
	// -1:异常态 1：初始态 2：备机状态 3：主机状态
	switch global.SelfFlag {
	case -1: //异常态
	case 1: //初始态
		for i, item := range global.SingletonNodeInfo.Clusters {
			code, info, err := checkUrl(item.Address)
			if err != nil {
				log.Warn("checkUrl error:" + err.Error())
				return err
			}
			if code == "3" {
				if err := cluster.SynchronyData(info.Address); err != nil {
					log.Warn("初始态->异常态 数据同步不成功:" + err.Error())
					global.SelfFlag = -1
					return err
				} else {
					log.Warn("初始态->备机状态")
					global.SelfFlag = 2
					global.MasterUrl = info.Address
				}
				break
			}
			if i+1 == len(global.SingletonNodeInfo.Clusters) {
				global.MasterUrl = global.CuCluster.Address
				global.SelfFlag = 3
				log.Warn("初始态->主机状态")
			}
		}
	case 2: //备机状态
		code, _, err := checkUrl(global.MasterUrl)
		if err != nil {
			log.Warn("check master Url error:" + err.Error())
			return err
		}
		if code == "3" {
			return nil
		} else {
			for i, item := range global.SingletonNodeInfo.Clusters {
				code, info, err := checkUrl(item.Address)
				if err != nil {
					log.Warn("check Url error:" + err.Error())

				} else if code == "3" {
					log.Warn("备机状态->备机状态")
					global.SelfFlag = 2
					if global.MasterUrl != info.Address {
						log.Warn("主机已切换: " + global.MasterUrl + "->" + info.Address + "")
						global.MasterUrl = info.Address
					}
					return nil
				}
				if i+1 == len(global.SingletonNodeInfo.Clusters) {
					log.Warn("备机状态->主机状态")
					global.SelfFlag = 3
					global.MasterUrl = global.CuCluster.Address
				}
			}
		}
	case 3: //主机状态
		for _, item := range global.SingletonNodeInfo.Clusters {
			code, info, err := checkUrl(item.Address)
			if err != nil {
				log.Warn("check Url error:" + err.Error())

			} else if code == "3" {
				if info.Level > global.CuCluster.Level {
					log.Warn("当前master【address:" + global.CuCluster.Address + " level:" + strconv.Itoa(global.CuCluster.Level) + " name:" + global.CuCluster.Name + "】 发现 另外 master【address:" + info.Address + " level:" + strconv.Itoa(info.Level) + " name:" + info.Name + "】")
					global.SelfFlag = 2
					global.MasterUrl = info.Address
					log.Warn("主机状态->备机状态 (遇到权重更好主机)")
				}
			}
		}
	default:
		global.SelfFlag = -1
		log.Error("当前机器状态" + strconv.Itoa(global.SelfFlag) + "异常 停止检测")
		return nil
	}
	return nil
}

func checkUrl(url string) (string, *models.Cluster, error) {
	request, err := http.NewRequest("GET", "http://"+url+"/api/IsMaster/", nil)
	if err != nil {
		return "", nil, err
	}
	response, err := new(http.Client).Do(request)
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
	var backJsonObj = cluster.ClusterBackObj{Data: &dataMsg}
	err = json.Unmarshal([]byte(bodyStr), &backJsonObj)
	if err != nil {
		return "", nil, err
	}
	var otherMaster *models.Cluster
	err = json.Unmarshal(dataMsg, otherMaster)
	if err != nil {
		return "", nil, err
	}
	return backJsonObj.Code, otherMaster, err
}
