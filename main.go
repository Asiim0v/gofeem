package main

import (
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/zserge/lorca"
)

func main() {
	// Gin 协程
	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		router.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		router.Run(":8080")
	}()

	var ui lorca.UI
	ui, _ = lorca.New("http://localhost:8080/ping", "", 800, 600)

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)

	select {
	case <-chSignal:
	case <-ui.Done():
	}
	ui.Close()
}
