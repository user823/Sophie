// Code generated by Kitex v0.8.0. DO NOT EDIT.
package genservice

import (
	server "github.com/cloudwego/kitex/server"
	v1 "github.com/user823/Sophie/api/thrift/gen/v1"
)

// NewServer creates a server.Server with the given handler and options.
func NewServer(handler v1.GenService, opts ...server.Option) server.Server {
	var options []server.Option

	options = append(options, opts...)

	svr := server.NewServer(options...)
	if err := svr.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	return svr
}
