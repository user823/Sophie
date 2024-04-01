package rpc

import (
	"fmt"
	"github.com/cloudwego/kitex/client"
	"sync"
)

type RPCClient struct {
	*WorkerClient
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
		Remoting = &RPCClient{
			WorkerClient: newWorkerClient(),
		}

		Remoting.WorkerClient.initRPC(generalOpts)
	})
}
