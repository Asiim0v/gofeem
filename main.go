package main

import (
	"embed"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zserge/lorca"
)

//go:embed frontend/dist/*
var FS embed.FS

// POST /api/v1/tests
// 1. 获取 go 执行文件所在目录
// 2. 在该目录创建 uploads 目录
// 3. 将文本保存为一个文件
// 4. 返回该文件的下载路径
func TextsController(c *gin.Context) {
	var json struct {
		Raw string
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		exe, err := os.Executable() // 获取当前执行文件的路径
		if err != nil {
			log.Fatal(err)
		}
		dir := filepath.Dir(exe) // 获取当前执行文件所在目录
		if err != nil {
			log.Fatal(err)
		}
		filename := uuid.New().String()          // 创建一个文件名
		uploads := filepath.Join(dir, "uploads") // 拼接 uploads 绝对路径
		err = os.MkdirAll(uploads, os.ModePerm)  // 创建 uploads 目录
		if err != nil {
			log.Fatal(err)
		}
		fullpath := path.Join("uploads", filename+".txt")
		err = ioutil.WriteFile(filepath.Join(dir, fullpath), []byte(json.Raw), 0644) // 将 json.Raw 写入文件
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"url": "/" + fullpath}) // 返回文件绝对路径(不含目录)
	}
}

func main() {
	// Gin 协程
	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()

		// // websocket, Phone 端上传时通知 PC 端
		// router.GET("/ws", func(c *gin.Context) {
		// 	ws.HttpController(c, hub)
		// })
		// // 下载接口, uploads 为所有上次的文件, path 为目标路径
		// router.GET("/uploads/:path", controllers.UploadsController)
		// // address 获取当前局域网 IP
		// router.GET("/api/v1/addresses", controllers.AddressesController)
		// // qrcode 局域网 IP 转为二维码
		// router.GET("/api/v1/qrcodes", controllers.QrcodesController)
		// // files 上传文件
		// router.POST("/api/v1/files", controllers.FilesController)
		// texts 上传文本
		router.POST("/api/v1/texts", TextsController)
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
		router.Run(":8080")
	}()

	var ui lorca.UI
	ui, _ = lorca.New("http://localhost:8080/static/index.html", "", 800, 600)

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)

	select {
	case <-chSignal:
	case <-ui.Done():
	}
	ui.Close()
}
