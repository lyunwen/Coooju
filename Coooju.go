package main

import (
	"./api"
	"./cluster"
	"./common/log"
	"./sockects"
	"./timer"
	"github.com/gin-gonic/gin"
)

func main() {
	cluster.Init()
	router := gin.Default()                                                 //api路由       //主同步备接口
	router.Group("/api/cluster/getData").GET("/", api.GetData)              //备机找主机
	router.Group("/api/cluster/syncData").GET("/", api.SyncData)            //备机找主机
	router.Group("/api/cluster/getCusterInfo").GET("/", api.GetClusterInfo) //备机找主机
	router.Group("/api/cluster/getNodeInfo").GET("/", api.GetNodeInfo)      //备机找主机
	//web socket 路由
	router.GET("/ws", func(c *gin.Context) { sockects.WebSocketHandler(c.Writer, c.Request) })
	//html页面路由
	//router.LoadHTMLGlob("view/*")
	//router.Group("/view/").GET("/:name", func(c *gin.Context) { c.HTML(http.StatusOK, c.Param("name")+".html", gin.H{}) })
	//静态文件路由
	router.Static("/view", "./view")
	router.Static("/wwwroot", "./wwwroot")

	log.Warn("==============================启动：" + cluster.CurrentData.Address + "==============================")
	go timer.Load()
	_ = router.Run(cluster.CurrentData.Address)
}
