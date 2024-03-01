package utils

import (
	"net"
	"os"
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func IsValidIP(ip string) bool {
	// 使用net.ParseIP来解析IP地址
	addr := net.ParseIP(ip)

	// 如果解析成功，则是合法的IP地址
	return addr != nil
}
