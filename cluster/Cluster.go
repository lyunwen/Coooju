package cluster

import (
	"../models"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type BackObj struct {
	Code string
	Msg  string
	Data interface{}
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
