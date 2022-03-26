package controllers

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/addresses
// 1. 获取电脑在各个局域网的 IP 地址
// 2、转为 JSON 写入 HTTP 响应
func AddressesController(c *gin.Context) {
	addrs, _ := net.InterfaceAddrs() // 获取 PC 在各个局域网的 IP 地址， 用户选择 Phone 所在的对应 IP 地址
	var result []string
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				result = append(result, ipnet.IP.String())
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"addresses": result})
}
