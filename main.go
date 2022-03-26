package main

import (
	"os"
	"os/signal"

	"github.com/Asiim0v/gofeem/server"
	"github.com/zserge/lorca"
)

func main() {

	// Gin 协程
	go server.Run()

	// lorca 启动 Chrome
	var ui lorca.UI
	ui, _ = lorca.New("http://localhost:"+server.Port+"/static/index.html", "", 800, 600)

	// 监听中断信号
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)

	// 等待中断信号
	select {
	case <-chSignal:
	case <-ui.Done():
	}
	ui.Close()
}
