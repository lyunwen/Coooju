package models

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"sync"
)

type Data struct {
	Version     string //主节点标识+版本号
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

func (data *Data) GetData() *Data {
	dataName := "data.json"
	dataJsonByte, err := ioutil.ReadFile(dataName)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes.TrimPrefix(dataJsonByte, []byte("\xef\xbb\xbf")), &data)
	if err != nil {
		panic(err)
	}
	return data
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

func (data *Data) SetData() error {
	mutex.Lock()
	dataJsonByte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	dataJsonStr := string(dataJsonByte)
	err = ioutil.WriteFile("data.json", []byte(dataJsonStr), 0644)
	if err != nil {
		return err
	}
	mutex.Unlock()
	return nil
}

func (data *Data) CopyData(fileName string) {
	mutex.Lock()
	dataJsonByte, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	dataJsonStr := string(dataJsonByte)
	err = ioutil.WriteFile(fileName+".json", []byte(dataJsonStr), 0644)
	if err != nil {
		panic(err)
	}
	mutex.Unlock()
}
