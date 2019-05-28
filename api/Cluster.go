package api

import (
	"../global"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
)

//func IsMaster(c *gin.Context) {
//	c.JSON(http.StatusOK, gin.H{"code": strconv.Itoa(global.SelfFlag), "msg": "", "data": global.CuCluster})
//}

//func GetMasterAddress(c *gin.Context) {
//	c.JSON(http.StatusOK, gin.H{"code": strconv.Itoa(global.SelfFlag), "msg": "", "data": global.MasterUrl})
//}

func GetClusterInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "", "data": global.ClusterData})
}

func GetNodeInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "", "data": global.CurrentData})
}

var getVotesMutex sync.Mutex

func SetVotes(c *gin.Context) {
	term, err := strconv.Atoi(c.Query("term"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "99", "msg": "term error", "data": ""})
		return
	}
	getVotesMutex.Lock()
	defer func() {
		getVotesMutex.Unlock()
	}()
	if term < global.CurrentData.VotedTerm {
		c.JSON(http.StatusOK, gin.H{"code": "99", "msg": "lower term", "data": ""})
		return
	}
	if term == global.CurrentData.VotedTerm {
		if global.CurrentData.VotedState == global.VotedState_UnDo {
			global.CurrentData.VotedState = global.VotedState_Done
			c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "data": ""})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"code": "99", "msg": "has voted", "data": ""})
			return
		}
	}
	if term > global.CurrentData.VotedTerm {
		global.CurrentData.VotedState = global.VotedState_Done
		global.CurrentData.VotedTerm = term
		c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "data": ""})
		return
	}
}

func SyncData(c *gin.Context) {
	//err := cluster.SynchronyData(global.MasterUrl)
	//if err == nil {
	//	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "", "data": nil})
	//} else {
	//	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": err.Error(), "data": nil})
	//}
}
