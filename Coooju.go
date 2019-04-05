package main

import (
	"./api"
	"./cluster"
	"./data"
	"./sockects"
	"./timer"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	beforeStart()

	router := gin.Default()                                                      //api路由
	router.Group("/api/SynchronyNodeData").GET("/", api.SynchronyNodeData)       //主同步备接口
	router.Group("/api/IsMaster").GET("/", api.IsMaster)                         //备机找主机
	router.Group("/api/cluster/getData").GET("/", api.GetClusterData)            //备机找主机
	router.Group("/api/cluster/getMasterAddress").GET("/", api.GetMasterAddress) //获取主机地址
	//router.Group("/api/getAwards").GET("/", api.GetAwards)
	//router.Group("/api/initData").GET("/", api.InitData)
	//router.Group("/api/getNextAction").GET("/", api.GetNextAction)
	//router.Group("/api/ndraw").GET("/", api.NDraw)
	//router.Group("/api/exdraw").GET("/", api.ExDraw)
	//router.Group("/api/pooldraw").GET("/", api.PoolDraw)
	//router.Group("/api/addMoney").GET("/", api.AddPoolMoney)
	//router.Group("/api/initSystem").GET("/", api.InitSystem)
	//web socket 路由
	router.GET("/ws", func(c *gin.Context) { sockects.WebSocketHandler(c.Writer, c.Request) })

	//html页面路由
	router.LoadHTMLGlob("view/*")
	router.Group("/view/").GET("/:name", func(c *gin.Context) { c.HTML(http.StatusOK, c.Param("name")+".html", gin.H{}) })

	////静态文件路由
	router.Static("/wwwroot", "./wwwroot")

	url, err := cluster.GetAvailablePortAddress()
	if err != nil {
		fmt.Println("start error:" + err.Error())
		return
	}

	go afterStart()
	global.LocalUrl = url
	_ = router.Run(url)
}

func beforeStart() {
	if err := data.Load(); err != nil {
		panic(err)
	}
}

func afterStart() {
	timer.Load()
}
