package api

import (
	"../cluster"
	"../global"
	"../models"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
)

func SynchronyNodeData(c *gin.Context) {
	var dataStr = c.Query("data")
	dataStr, err := url.PathUnescape(dataStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "99", "msg": "data Unescape error", "data": nil})
		return
	}
	tranDataObj, err := new(models.Data).GetDataFromJsonStr(dataStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "99", "msg": "data explain error", "data": nil})
		return
	}
	result := tranDataObj.SetData()
	switch result {
	case "ok":
		global.SingletonNodeInfo = tranDataObj
		c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "data": nil})
		return
	case "smaller":
		c.JSON(http.StatusOK, gin.H{"code": "1", "msg": "smaller version", "data": nil})
		return
	case "equal":
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "equal version", "data": nil})
		return
	default:
		c.JSON(http.StatusOK, gin.H{"code": "3", "msg": "other error", "data": nil})
		return
	}
}

func IsMaster(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": strconv.Itoa(global.SelfFlag), "msg": "", "data": nil})
}

func GetMasterAddress(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": strconv.Itoa(global.SelfFlag), "msg": "", "data": global.MasterUrl})
}

func GetData(c *gin.Context) {
	data := new(models.Data).GetData()
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "", "data": data})
}

func SyncData(c *gin.Context) {
	err := cluster.SynchronyData()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "", "data": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": "0", "msg": err.Error(), "data": nil})
	}
}
