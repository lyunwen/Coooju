package models

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
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

func (data *Data) GetVersionInfo() (string, int) {
	var infos = strings.Split(data.Version, "-")
	version, err := strconv.Atoi(infos[1])
	if err != nil {
		panic(err)
	}
	return infos[0], version
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

func (data *Data) SetData() string {
	mutex.Lock()
	preData := new(Data).GetData()
	if preData.Version < data.Version {
		return "smaller"
	} else if preData.Version == data.Version {
		return "equal"
	} else {
		dataJsonByte, err := json.Marshal(data)
		if err != nil {
			return ""
		}
		dataJsonStr := string(dataJsonByte)
		err = ioutil.WriteFile("data.json", []byte(dataJsonStr), 0644)
		if err != nil {
			panic(err)
		}
		mutex.Unlock()
		return "ok"
	}
}
