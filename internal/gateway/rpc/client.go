package rpc

import (
	"fmt"
	"github.com/cloudwego/kitex/client"
	"sync"
)

type RPCClient struct {
	*SystemClient
	*JobClient
	*FileClient
	*GenClient
}

var (
	ErrRPC = fmt.Errorf("系统内部错误，请重试")
)

var (
	once     sync.Once
	Remoting *RPCClient
)

func Init(generalOpts []client.Option) {
	once.Do(func() {
		// ---使用cfg 初始化各个client组件---
		Remoting = &RPCClient{
			SystemClient: newSystemClient(),
			JobClient:    newJobClient(),
			FileClient:   newFileClient(),
			GenClient:    newGenClient(),
		}

		Remoting.SystemClient.initRPC(generalOpts)
		Remoting.JobClient.initRPC(generalOpts)
		Remoting.FileClient.initRPC(generalOpts)
		Remoting.GenClient.initRPC(generalOpts)
	})
}
