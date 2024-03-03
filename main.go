package main

import (
	v1 "github.com/user823/Sophie/api/thrift/system/v1/systemservice"
	"log"
)

func main() {
	svr := v1.NewServer(new(SystemServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
