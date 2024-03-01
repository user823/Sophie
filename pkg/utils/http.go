package utils

import (
	"github.com/cloudwego/hertz/pkg/app"
	"strings"
)

func GetClientIP(c *app.RequestContext) string {
	clientIp := c.ClientIP()
	if clientIp != "127.0.0.1" {
		return clientIp
	}

	remoteAddr := c.RemoteAddr().String()
	remoteIpAndPort := strings.Split(remoteAddr, ":")
	return remoteIpAndPort[0]
}
