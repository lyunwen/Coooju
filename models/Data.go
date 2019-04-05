package models

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"sync"
)

//主机使用或切主使用
var dataSingleton *Data

type Data struct {
	Version     int
	Description string
	Url         string
	Services    []Service
	Clusters    []Cluster
}

type Cluster struct {
	Level   int
	Name    string
	Address string
}

type Service struct {
	Name  string
	Url   string
	Nodes []Node
}

type Node struct {
	Name    string
	Address string
}

func (data *Data) GetData() (*Data, error) {
	dataName := "data.json"
	dataJsonByte, err := ioutil.ReadFile(dataName)
	if err != nil {
		return data, err
	}
	dataJsonStr := string(dataJsonByte)
	err = json.Unmarshal([]byte(dataJsonStr), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (data *Data) GetDataFromJsonStr(dataJsonStr string) (*Data, error) {
	err := json.Unmarshal([]byte(dataJsonStr), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//线程安全
var mutex sync.Mutex

func (data *Data) SetData() (string, error) {
	mutex.Lock()
	preData, err := new(Data).GetData()
	if err != nil {
		return "", err
	}
	if preData.Version < data.Version {
		return "", errors.New("smaller version error")
	} else if preData.Version == data.Version {
		return "equal", errors.New("equal version error")
	} else {
		dataJsonByte, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		dataJsonStr := string(dataJsonByte)
		err = ioutil.WriteFile("data.json", []byte(dataJsonStr), 0644)
		if err != nil {
			return "", err
		}
		mutex.Unlock()
		return "ok", nil
	}
}
