package main

import (
	"./api"
	"./cluster"
	"./common/log"
	"./data"
	"./sockects"
	"./timer"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	data.Load()
	router := gin.Default() //api路由       //主同步备接口
	//router.Group("/api/IsMaster").GET("/", api.IsMaster)                         //备机找主机
	router.Group("/api/setVotes").GET("/", api.SetVotes)                    //投票
	router.Group("/api/cluster/getCusterInfo").GET("/", api.GetClusterInfo) //备机找主机
	router.Group("/api/cluster/getNodeInfo").GET("/", api.GetNodeInfo)      //备机找主机
	//router.Group("/api/cluster/syncData").GET("/", api.SyncData)                 //备机找主机
	//router.Group("/api/cluster/getMasterAddress").GET("/", api.GetMasterAddress) //获取主机地址
	//web socket 路由
	router.GET("/ws", func(c *gin.Context) { sockects.WebSocketHandler(c.Writer, c.Request) })
	//html页面路由
	//router.LoadHTMLGlob("view/*")
	//router.Group("/view/").GET("/:name", func(c *gin.Context) { c.HTML(http.StatusOK, c.Param("name")+".html", gin.H{}) })
	//静态文件路由
	router.Static("/view", "./view")
	router.Static("/wwwroot", "./wwwroot")

	var url, err = cluster.GetAvailablePortAddress()
	if err != nil {
		fmt.Println("start error:" + err.Error())
		return
	}
	log.Warn("==============================启动：" + url + "==============================")
	go timer.Load()
	_ = router.Run(url)
}
