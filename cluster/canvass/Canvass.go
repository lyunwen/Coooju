package canvass

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

//拉票变量
var (
	HasVoted  bool
	ThisMutex sync.Mutex
)

//拉票
func Canvass(c *gin.Context) {
	if HasVoted == false {
		ThisMutex.Lock()
		if HasVoted == false {
			HasVoted = true
			c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "", "data": "1"})
		}
		ThisMutex.Unlock()
	}
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "", "data": "2"})
}
