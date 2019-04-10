package main

import (
	"./api"
	"./cluster"
	"./common/log"
	"./data"
	"./global"
	"./sockects"
	"./timer"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	log.Default("启动")
	data.Load()
	router := gin.Default()                                                      //api路由       //主同步备接口
	router.Group("/api/IsMaster").GET("/", api.IsMaster)                         //备机找主机
	router.Group("/api/cluster/getData").GET("/", api.GetData)                   //备机找主机
	router.Group("/api/cluster/syncData").GET("/", api.SyncData)                 //备机找主机
	router.Group("/api/cluster/getMasterAddress").GET("/", api.GetMasterAddress) //获取主机地址
	//web socket 路由
	router.GET("/ws", func(c *gin.Context) { sockects.WebSocketHandler(c.Writer, c.Request) })
	//html页面路由
	router.LoadHTMLGlob("view/*")
	router.Group("/view/").GET("/:name", func(c *gin.Context) { c.HTML(http.StatusOK, c.Param("name")+".html", gin.H{}) })
	//静态文件路由
	router.Static("/wwwroot", "./wwwroot")

	url, err := cluster.GetAvailablePortAddress()
	if err != nil {
		fmt.Println("start error:" + err.Error())
		return
	}
	global.LocalUrl = url
	go timer.Load()
	_ = router.Run(url)
}
