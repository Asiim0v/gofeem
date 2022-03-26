package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/Asiim0v/gofeem/server/controllers"
	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist/*
var FS embed.FS

var Port = "27149"

func Run() {

	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	// // websocket, Phone 端上传时通知 PC 端
	// router.GET("/ws", func(c *gin.Context) {
	// 	ws.HttpController(c, hub)
	// })
	// 下载接口, uploads 为所有上传的文件, path 为目标路径
	router.GET("/uploads/:path", controllers.UploadsController)
	// address 获取当前局域网 IP
	router.GET("/api/v1/addresses", controllers.AddressesController)
	// qrcode 局域网 IP 转为二维码
	router.GET("/api/v1/qrcodes", controllers.QrcodesController)
	// files 上传文件
	router.POST("/api/v1/files", controllers.FilesController)
	// texts 上传文本
	router.POST("/api/v1/texts", controllers.TextsController)
	// 添加静态路由, 所有 static 开头的路由都会自动读取 dist 下的对应文件
	staticFiles, _ := fs.Sub(FS, "frontend/dist")
	router.StaticFS("/static", http.FS(staticFiles))
	// 如果用户输入的路径不能匹配 dist 下的文件
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// 如果输入的是 static 开头的路径, 统一渲染成 index.html
		if strings.HasPrefix(path, "/static/") {
			reader, err := staticFiles.Open("index.html")
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()
			stat, err := reader.Stat()
			if err != nil {
				log.Fatal(err)
			}
			c.DataFromReader(http.StatusOK, stat.Size(), "text/html;charset=utf-8", reader, nil)
		} else {
			c.Status(http.StatusNotFound)
		}
	})
	router.Run(":" + Port)
}
