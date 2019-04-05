package api

import (
	"../global"
	"../models"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
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
	result, err := tranDataObj.SetData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "99", "msg": "set data error", "data": nil})
		return
	}
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
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "", "data": global.SelfFlag})
	return
}

func GetClusterData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "", "data": new(models.Data).GetData()})
	return
}
