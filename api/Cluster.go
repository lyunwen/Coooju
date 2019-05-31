package api

import (
	"../global"
	"../models/clusterState"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
)

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
		if global.CurrentData.ClusterState == clusterState.Leader {
			c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "lower term", "data": "found leader"})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "lower term", "data": "lower term"})
			return
		}
	}
	if term == global.CurrentData.VotedTerm {
		if global.CurrentData.ClusterState == clusterState.Leader {
			c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "lower term", "data": "found leader"})
			return
		} else if global.CurrentData.VotedState == global.VotedState_UnDo {
			global.CurrentData.VotedState = global.VotedState_Done
			c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "data": "ok"})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "has voted", "data": "has voted"})
			return
		}
	}
	if term > global.CurrentData.VotedTerm {
		global.CurrentData.VotedState = global.VotedState_Done
		global.CurrentData.VotedTerm = term
		c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "data": "ok"})
		return
	}
}
