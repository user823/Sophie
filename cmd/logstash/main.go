package main

import "github.com/user823/Sophie/internal/logstash"

func main() {
	logstash.NewApp("sophie-logstash").Run()
}
