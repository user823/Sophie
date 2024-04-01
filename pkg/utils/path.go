package utils

import (
	"net"
)

func IsValidIP(ip string) bool {
	// 使用net.ParseIP来解析IP地址
	addr := net.ParseIP(ip)

	// 如果解析成功，则是合法的IP地址
	return addr != nil
}
