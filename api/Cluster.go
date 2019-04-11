package api

import (
	"../cluster"
	"../global"
	"../models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func IsMaster(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": strconv.Itoa(global.SelfFlag), "msg": "", "data": global.CuCluster})
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
