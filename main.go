package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zserge/lorca"
)

//go:embed frontend/dist/*
var FS embed.FS

func main() {
	// Gin 协程
	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
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
