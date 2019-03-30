package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var (
	voteLevel  int
	hasVoted   bool
	isMaster   bool
	votedCount int
)

func InitVotedData() {
	voteLevel = 1
	hasVoted = false
	isMaster = false
	votedCount = 0
}

func Vote(c *gin.Context) {
	beVoteLevel, err := strconv.Atoi(c.Query("level"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "99", "msg": "level error"})
		return
	}
	if voteLevel >= beVoteLevel {
		return
	} else {
		return
	}
}
