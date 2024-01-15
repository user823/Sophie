package main

import (
	v1 "github.com/user823/Sophie/api/thrift/file/v1/fileservice"
	"log"
)

func main() {
	svr := v1.NewServer(new(FileServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
