package main

import "github.com/user823/Sophie/internal/gateway"

// @title sophie-gateway
// @version 1.0
// @license.name Apache 2.0
func main() {
	gateway.NewApp("sophie-gateway").Run()
}
