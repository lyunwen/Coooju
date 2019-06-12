package api

import (
	"../cluster"
	"../models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetClusterInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "", "data": cluster.ShareData})
}

func GetNodeInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "", "data": cluster.CurrentNodeInfo{
		MasterAddress: cluster.OwnData.MasterAddress,
		ClusterState:  cluster.OwnData.ClusterState,
		Address:       cluster.OwnData.Address,
		Name:          cluster.OwnData.GetName(),
		Level:         cluster.OwnData.GetLevel(),
	}})
}

func GetData(c *gin.Context) {
	data := new(models.Data).GetData()
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "", "data": data})
}

func SyncData(c *gin.Context) {
	err := cluster.SynchronyData(cluster.OwnData.MasterAddress)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "", "data": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": "0", "msg": err.Error(), "data": nil})
	}
}
