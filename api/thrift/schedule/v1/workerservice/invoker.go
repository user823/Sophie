// Code generated by Kitex v0.8.0. DO NOT EDIT.

package workerservice

import (
	server "github.com/cloudwego/kitex/server"
	v1 "github.com/user823/Sophie/api/thrift/schedule/v1"
)

// NewInvoker creates a server.Invoker with the given handler and options.
func NewInvoker(handler v1.WorkerService, opts ...server.Option) server.Invoker {
	var options []server.Option

	options = append(options, opts...)

	s := server.NewInvoker(options...)
	if err := s.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	if err := s.Init(); err != nil {
		panic(err)
	}
	return s
}
